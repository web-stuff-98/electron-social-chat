package socketserver

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/web-stuff-98/electron-social-chat/pkg/db"
	"github.com/web-stuff-98/electron-social-chat/pkg/socketmodels"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*--------------- SOCKET SERVER STRUCT ---------------*/
type SocketServer struct {
	Connections                 Connections
	Subscriptions               Subscriptions
	ConnectionSubscriptionCount ConnectionsSubscriptionCount

	RegisterConn   chan ConnectionInfo
	UnregisterConn chan ConnectionInfo

	RegisterSubscriptionConn   chan SubscriptionConnectionInfo
	UnregisterSubscriptionConn chan SubscriptionConnectionInfo

	SendDataToSubscription           chan SubscriptionDataMessage
	SendDataToSubscriptionExclusive  chan ExclusiveSubscriptionDataMessage
	SendDataToSubscriptions          chan SubscriptionDataMessageMulti
	SendDataToSubscriptionsExclusive chan ExclusiveSubscriptionDataMessageMulti
	RemoveUserFromSubscription       chan RemoveUserFromSubscription

	// websocket Write/Read must be done from 1 goroutine. Queue all of them to be executed in a loop.
	MessageSendQueue chan QueuedMessage

	SendDataToUser chan UserDataMessage

	DestroySubscription chan string
}

/*--------------- MUTEX PROTECTED MAPS ---------------*/
type Connections struct {
	data  map[*websocket.Conn]primitive.ObjectID
	mutex sync.Mutex
}
type Subscriptions struct {
	data  map[string]map[*websocket.Conn]primitive.ObjectID
	mutex sync.Mutex
}

/*--------------- OTHER STRUCTS ---------------*/
type ConnectionInfo struct {
	Conn   *websocket.Conn
	Uid    primitive.ObjectID
	Online bool
}
type ConnectionsSubscriptionCount struct {
	data  map[*websocket.Conn]uint8 //Max subscriptions is 128... nice number half max uint8
	mutex sync.Mutex
}
type SubscriptionConnectionInfo struct {
	Name string
	Uid  primitive.ObjectID
	Conn *websocket.Conn
}
type SubscriptionDataMessage struct {
	Name string
	Data []byte
}
type ExclusiveSubscriptionDataMessage struct {
	Name    string
	Data    []byte
	Exclude map[primitive.ObjectID]bool
}
type SubscriptionDataMessageMulti struct {
	Names []string
	Data  []byte
}
type ExclusiveSubscriptionDataMessageMulti struct {
	Names   []string
	Data    []byte
	Exclude map[primitive.ObjectID]bool
}
type QueuedMessage struct {
	Conn *websocket.Conn
	Data []byte
}
type UserDataMessage struct {
	Uid  primitive.ObjectID
	Data interface{}
	Type string
}
type RemoveUserFromSubscription struct {
	Name string
	Uid  primitive.ObjectID
}

func Init(colls *db.Collections) (*SocketServer, error) {
	socketServer := &SocketServer{
		Connections: Connections{
			data: make(map[*websocket.Conn]primitive.ObjectID),
		},
		Subscriptions: Subscriptions{
			data: make(map[string]map[*websocket.Conn]primitive.ObjectID),
		},
		ConnectionSubscriptionCount: ConnectionsSubscriptionCount{
			data: make(map[*websocket.Conn]uint8),
		},
		RegisterConn:   make(chan ConnectionInfo),
		UnregisterConn: make(chan ConnectionInfo),

		RegisterSubscriptionConn:   make(chan SubscriptionConnectionInfo),
		UnregisterSubscriptionConn: make(chan SubscriptionConnectionInfo),

		SendDataToSubscription:           make(chan SubscriptionDataMessage),
		SendDataToSubscriptionExclusive:  make(chan ExclusiveSubscriptionDataMessage),
		SendDataToSubscriptions:          make(chan SubscriptionDataMessageMulti),
		SendDataToSubscriptionsExclusive: make(chan ExclusiveSubscriptionDataMessageMulti),
		RemoveUserFromSubscription:       make(chan RemoveUserFromSubscription),

		MessageSendQueue: make(chan QueuedMessage),

		SendDataToUser: make(chan UserDataMessage),

		DestroySubscription: make(chan string),
	}
	RunServer(socketServer, colls)
	return socketServer, nil
}

func RunServer(socketServer *SocketServer, colls *db.Collections) {
	/* ----- Connection registration ----- */
	go func() {
		for {
			defer func() {
				r := recover()
				if r != nil {
					log.Println("Recovered from panic in WS registration :", r)
				}
			}()
			connData := <-socketServer.RegisterConn
			if connData.Conn != nil {
				socketServer.Connections.mutex.Lock()
				socketServer.Connections.data[connData.Conn] = connData.Uid
				socketServer.Connections.mutex.Unlock()
				outBytes, err := json.Marshal(socketmodels.OutChangeMessage{
					Type:   "CHANGE",
					Method: "UPDATE",
					Data:   `{"ID":"` + connData.Uid.Hex() + `"` + `,"online":true}`,
					Entity: "USER",
				})
				if err == nil {
					socketServer.SendDataToSubscription <- SubscriptionDataMessage{
						Name: "user=" + connData.Uid.Hex(),
						Data: outBytes,
					}
				}
			}
		}
	}()
	/* ----- Disconnect registration ----- */
	go func() {
		for {
			defer func() {
				r := recover()
				if r != nil {
					log.Println("Recovered from panic in WS deregistration :", r)
				}
			}()
			connData := <-socketServer.UnregisterConn
			socketServer.Connections.mutex.Lock()
			socketServer.Subscriptions.mutex.Lock()
			for conn := range socketServer.Connections.data {
				if conn == connData.Conn {
					delete(socketServer.Connections.data, conn)
					for _, r := range socketServer.Subscriptions.data {
						for c := range r {
							if c == connData.Conn {
								delete(r, c)
								break
							}
						}
					}
					break
				}
			}
			if connData.Uid != primitive.NilObjectID {
				outBytes, err := json.Marshal(socketmodels.OutChangeMessage{
					Type:   "CHANGE",
					Method: "UPDATE",
					Data:   `{"ID":"` + connData.Uid.Hex() + `"` + `,"online":false}`,
					Entity: "USER",
				})
				if err == nil {
					socketServer.SendDataToSubscription <- SubscriptionDataMessage{
						Name: "user=" + connData.Uid.Hex(),
						Data: outBytes,
					}
				}
			}
			socketServer.Connections.mutex.Unlock()
			socketServer.Subscriptions.mutex.Unlock()
		}
	}()
	/* ----- Send messages in queue ----- */
	go func() {
		for {
			defer func() {
				r := recover()
				if r != nil {
					log.Println("Recovered from panic in queued socket messages :", r)
				}
			}()
			data := <-socketServer.MessageSendQueue
			data.Conn.WriteMessage(websocket.TextMessage, data.Data)
		}
	}()
	/* ----- Subscription connection registration (also check the authorization if subscription requires it) ----- */
	go func() {
		for {
			defer func() {
				r := recover()
				if r != nil {
					log.Println("Recovered from panic in subscription registration :", r)
				}
			}()
			connData := <-socketServer.RegisterSubscriptionConn
			allow := true
			// Make sure users cannot open too many subscriptions
			socketServer.ConnectionSubscriptionCount.mutex.Lock()
			socketServer.Subscriptions.mutex.Lock()
			count, countOk := socketServer.ConnectionSubscriptionCount.data[connData.Conn]
			if count >= 128 {
				allow = false
			}
			if connData.Conn != nil {
				// Passed all checks, add the connection to the subscription
				if allow {
					if socketServer.Subscriptions.data[connData.Name] == nil {
						socketServer.Subscriptions.data[connData.Name] = make(map[*websocket.Conn]primitive.ObjectID)
					}
					socketServer.Subscriptions.data[connData.Name][connData.Conn] = connData.Uid
					if countOk {
						socketServer.ConnectionSubscriptionCount.data[connData.Conn]++
					} else {
						socketServer.ConnectionSubscriptionCount.data[connData.Conn] = 1
					}
				}
				socketServer.Subscriptions.mutex.Unlock()
				socketServer.ConnectionSubscriptionCount.mutex.Unlock()
			}
		}
	}()
	/* ----- Subscription disconnect registration ----- */
	go func() {
		for {
			defer func() {
				r := recover()
				if r != nil {
					log.Println("Recovered from panic in subscription disconnect registration :", r)
				}
			}()
			connData := <-socketServer.UnregisterSubscriptionConn
			var err error
			if connData.Conn == nil {
				err = fmt.Errorf("Connection was nil")
			}
			if err != nil {
				socketServer.Subscriptions.mutex.Lock()
				if _, ok := socketServer.Subscriptions.data[connData.Name]; ok {
					delete(socketServer.Subscriptions.data[connData.Name], connData.Conn)
				}
				socketServer.Subscriptions.mutex.Unlock()
				socketServer.ConnectionSubscriptionCount.mutex.Lock()
				if _, ok := socketServer.ConnectionSubscriptionCount.data[connData.Conn]; ok {
					socketServer.ConnectionSubscriptionCount.data[connData.Conn]--
				}
				socketServer.ConnectionSubscriptionCount.mutex.Unlock()
			}
		}
	}()
	/* ----- Send data to subscription ----- */
	go func() {
		for {
			defer func() {
				r := recover()
				if r != nil {
					log.Println("Recovered from panic in subscription data channel :", r)
				}
			}()
			subsData := <-socketServer.SendDataToSubscription
			socketServer.Subscriptions.mutex.Lock()
			for k, s := range socketServer.Subscriptions.data {
				if k == subsData.Name {
					for conn := range s {
						socketServer.MessageSendQueue <- QueuedMessage{
							Conn: conn,
							Data: subsData.Data,
						}
					}
					break
				}
			}
			socketServer.Subscriptions.mutex.Unlock()
		}
	}()
	/* ----- Send data to subscription excluding uids ----- */
	go func() {
		for {
			defer func() {
				r := recover()
				if r != nil {
					log.Println("Recovered from panic in exclusive subscription data channel :", r)
				}
			}()
			subsData := <-socketServer.SendDataToSubscriptionExclusive
			socketServer.Subscriptions.mutex.Lock()
			for k, s := range socketServer.Subscriptions.data {
				if k == subsData.Name {
					for conn, oid := range s {
						if subsData.Exclude[oid] != true {
							socketServer.MessageSendQueue <- QueuedMessage{
								Conn: conn,
								Data: subsData.Data,
							}
						}
					}
					break
				}
			}
			socketServer.Subscriptions.mutex.Unlock()
		}
	}()
	/* ----- Send data to multiple subscriptions ----- */
	go func() {
		for {
			defer func() {
				r := recover()
				if r != nil {
					log.Println("Recovered from panic in subscription data channel :", r)
				}
			}()
			subsData := <-socketServer.SendDataToSubscriptions
			socketServer.Subscriptions.mutex.Lock()
			for _, v := range subsData.Names {
				for k, s := range socketServer.Subscriptions.data {
					if k == v {
						for conn := range s {
							socketServer.MessageSendQueue <- QueuedMessage{
								Conn: conn,
								Data: subsData.Data,
							}
						}
						break
					}
				}
			}
			socketServer.Subscriptions.mutex.Unlock()
		}
	}()
	/* ----- Send data to multiple subscriptions excluding uids ----- */
	go func() {
		for {
			defer func() {
				r := recover()
				if r != nil {
					log.Println("Recovered from panic in exclusive subscription data channel :", r)
				}
			}()
			subsData := <-socketServer.SendDataToSubscriptionsExclusive
			socketServer.Subscriptions.mutex.Lock()
			for _, v := range subsData.Names {
				for k, s := range socketServer.Subscriptions.data {
					if k == v {
						for conn, oid := range s {
							if subsData.Exclude[oid] != true {
								socketServer.MessageSendQueue <- QueuedMessage{
									Conn: conn,
									Data: subsData.Data,
								}
							}
						}
						break
					}
				}
			}
			socketServer.Subscriptions.mutex.Unlock()
		}
	}()
	/* ----- Send data to a specific user ----- */
	go func() {
		for {
			defer func() {
				r := recover()
				if r != nil {
					log.Println("Recovered from panic in send data to user channel :", r)
				}
			}()
			data := <-socketServer.SendDataToUser
			socketServer.Connections.mutex.Lock()
			for conn, uid := range socketServer.Connections.data {
				if data.Uid == uid {
					var m map[string]interface{}
					outBytesNoTypeKey, err := json.Marshal(data.Data)
					json.Unmarshal(outBytesNoTypeKey, &m)
					m["TYPE"] = data.Type
					outBytes, err := json.Marshal(m)
					if err == nil {
						socketServer.MessageSendQueue <- QueuedMessage{
							Conn: conn,
							Data: outBytes,
						}
					} else {
						log.Println("Error marshaling data to be sent to user :", err)
					}
					break
				}
			}
			socketServer.Connections.mutex.Unlock()
		}
	}()
	/* ----- Remove a user from subscription ----- */
	go func() {
		for {
			defer func() {
				r := recover()
				if r != nil {
					log.Println("Recovered from panic in remove user from subscription channel :", r)
				}
			}()
			data := <-socketServer.RemoveUserFromSubscription
			socketServer.Subscriptions.mutex.Lock()
			if subs, ok := socketServer.Subscriptions.data[data.Name]; ok {
				for c, oi := range subs {
					if oi == data.Uid {
						defer func() {
							socketServer.Subscriptions.mutex.Unlock()
						}()
						socketServer.Subscriptions.mutex.Lock()
						delete(socketServer.Subscriptions.data[data.Name], c)
						break
					}
				}
			}
			socketServer.Subscriptions.mutex.Unlock()
		}
	}()
	/* ----- Destroy subscription ----- */
	go func() {
		for {
			defer func() {
				r := recover()
				if r != nil {
					log.Println("Recovered from panic in destroy subscription channel :", r)
				}
			}()
			subsName := <-socketServer.DestroySubscription
			socketServer.Subscriptions.mutex.Lock()
			socketServer.ConnectionSubscriptionCount.mutex.Lock()
			for c := range socketServer.Subscriptions.data[subsName] {
				if _, ok := socketServer.ConnectionSubscriptionCount.data[c]; ok {
					socketServer.ConnectionSubscriptionCount.data[c]--
				}
			}
			delete(socketServer.Subscriptions.data, subsName)
			socketServer.Subscriptions.mutex.Unlock()
			socketServer.ConnectionSubscriptionCount.mutex.Unlock()
		}
	}()
}
