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

/* --------------- SOCKET SERVER STRUCT --------------- */
type SocketServer struct {
	Connections                 Connections
	Subscriptions               Subscriptions
	ConnectionSubscriptionCount ConnectionsSubscriptionCount

	RegisterConn                       chan ConnectionInfo
	UnregisterConn                     chan ConnectionInfo
	AttachmentServerRemoveUploaderChan chan primitive.ObjectID

	RegisterSubscriptionConn   chan SubscriptionConnectionInfo
	UnregisterSubscriptionConn chan SubscriptionConnectionInfo

	SendDataToSubscription           chan SubscriptionDataMessage
	SendDataToSubscriptionExclusive  chan ExclusiveSubscriptionDataMessage
	SendDataToSubscriptions          chan SubscriptionDataMessageMulti
	SendDataToSubscriptionsExclusive chan ExclusiveSubscriptionDataMessageMulti
	RemoveUserFromSubscription       chan RemoveUserFromSubscription
	DestroySubscription              chan string
	GetSubscriptionUids              chan GetSubscriptionUids

	// Calls that have not yet been answered
	CallsPending CallsPending
	// Channel for creating calls
	CallsPendingChan chan InCall
	// Channel for responding to calls
	ResponseToCallChan chan InCallResponse
	// Mutex protected map for active calls
	CallsActive CallsActive

	// Websocket Write/Read must be done from 1 goroutine. Queue all messages to be sent 1 by 1.
	// This is probably a bad way to do it... Should be a seperate goroutine for each channel,
	// it would increase the speed of messaging by sending out messages to different connections
	// seperately instead of having them all wait in one long queue... really tricky to implement
	// though, too dumb to do it
	MessageSendQueue chan QueuedMessage

	SendDataToUser  chan UserDataMessage
	SendDataToUsers chan UsersDataMessage
}

/* --------------- MUTEX PROTECTED MAPS --------------- */
type Connections struct {
	data  map[*websocket.Conn]primitive.ObjectID
	mutex sync.Mutex
}
type Subscriptions struct {
	data  map[string]map[*websocket.Conn]primitive.ObjectID
	mutex sync.Mutex
}
type ConnectionsSubscriptionCount struct {
	data  map[*websocket.Conn]uint8 //Max subscriptions is 128... nice number half max uint8
	mutex sync.Mutex
}
type CallsPending struct {
	// outer map is caller ID, inner map is the user that was called ID
	data  map[primitive.ObjectID]primitive.ObjectID
	mutex sync.Mutex
}
type InCall struct {
	Caller primitive.ObjectID
	Called primitive.ObjectID
}
type InCallResponse struct {
	Caller primitive.ObjectID
	Called primitive.ObjectID
	Accept bool
}
type CallsActive struct {
	// outer map is caller ID, inner map is the user that was called ID
	data  map[primitive.ObjectID]primitive.ObjectID
	mutex sync.Mutex
}

/* --------------- OTHER STRUCTS --------------- */
type ConnectionInfo struct {
	Conn   *websocket.Conn
	Uid    primitive.ObjectID
	Online bool
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
type UsersDataMessage struct {
	Uids map[primitive.ObjectID]struct{}
	Data interface{}
	Type string
}
type RemoveUserFromSubscription struct {
	Name string
	Uid  primitive.ObjectID
}

/* --------------- RECV CHAN STRUCTS --------------- */
type GetSubscriptionUids struct {
	RecvChan chan<- map[primitive.ObjectID]struct{}
	Name     string
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

		RegisterConn:                       make(chan ConnectionInfo),
		UnregisterConn:                     make(chan ConnectionInfo),
		AttachmentServerRemoveUploaderChan: make(chan primitive.ObjectID),

		RegisterSubscriptionConn:   make(chan SubscriptionConnectionInfo),
		UnregisterSubscriptionConn: make(chan SubscriptionConnectionInfo),

		SendDataToSubscription:           make(chan SubscriptionDataMessage),
		SendDataToSubscriptionExclusive:  make(chan ExclusiveSubscriptionDataMessage),
		SendDataToSubscriptions:          make(chan SubscriptionDataMessageMulti),
		SendDataToSubscriptionsExclusive: make(chan ExclusiveSubscriptionDataMessageMulti),
		RemoveUserFromSubscription:       make(chan RemoveUserFromSubscription),
		DestroySubscription:              make(chan string),
		GetSubscriptionUids:              make(chan GetSubscriptionUids),

		CallsPending: CallsPending{
			data: make(map[primitive.ObjectID]primitive.ObjectID),
		},
		CallsPendingChan:   make(chan InCall),
		ResponseToCallChan: make(chan InCallResponse),
		CallsActive: CallsActive{
			data: make(map[primitive.ObjectID]primitive.ObjectID),
		},

		MessageSendQueue: make(chan QueuedMessage),

		SendDataToUser:  make(chan UserDataMessage),
		SendDataToUsers: make(chan UsersDataMessage),
	}
	RunServer(socketServer, colls)
	return socketServer, nil
}

func RunServer(socketServer *SocketServer, colls *db.Collections) {
	/* ----- Connection registration ----- */
	go connectionRegistrationLoop(socketServer, colls)
	/* ----- Disconnect registration ----- */
	go disconnectRegistrationLoop(socketServer, colls)
	/* ----- Send messages in queue ----- */
	go messageQueueLoop(socketServer, colls)
	/* ----- Subscription connection registration (also check the authorization if subscription requires it) ----- */
	go subscriptionConnectionRegistrationLoop(socketServer, colls)
	/* ----- Subscription disconnect registration ----- */
	go subscriptionDisconnectRegistrationLoop(socketServer, colls)
	/* ----- Send data to subscription ----- */
	go sendSubscriptionDataLoop(socketServer, colls)
	/* ----- Send data to subscription excluding uids ----- */
	go sendSubscriptionDataExclusiveLoop(socketServer, colls)
	/* ----- Send data to multiple subscriptions ----- */
	go sendToMultipleSubscriptionsLoop(socketServer, colls)
	/* ----- Send data to multiple subscriptions excluding uids ----- */
	go sendToMultipleSubscriptionsExclusiveLoop(socketServer, colls)
	/* ----- Get uids of users using subscription ----- */
	go getSubscriptionUidsLoop(socketServer, colls)
	/* ----- Send data to a specific user ----- */
	go sendDataToUserLoop(socketServer, colls)
	/* ----- Send data to users ----- */
	go sendDataToUsersLoop(socketServer, colls)
	/* ----- Remove a user from subscription ----- */
	go removeUserFromSubscriptionLoop(socketServer, colls)
	/* ----- Destroy subscription ----- */
	go destroySubscriptionLoop(socketServer, colls)
	/* ----- Call pending loop ----- */
	go callPendingChanLoop(socketServer, colls)
	/* ----- Call response loop ----- */
	go callResponseChanLoop(socketServer, colls)
}

func connectionRegistrationLoop(socketServer *SocketServer, colls *db.Collections) {
	for {
		defer func() {
			r := recover()
			if r != nil {
				log.Println("Recovered from panic in WS registration :", r)
			}
			go connectionRegistrationLoop(socketServer, colls)
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
}

func disconnectRegistrationLoop(socketServer *SocketServer, colls *db.Collections) {
	for {
		defer func() {
			r := recover()
			if r != nil {
				log.Println("Recovered from panic in WS deregistration :", r)
			}
			go disconnectRegistrationLoop(socketServer, colls)
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
		socketServer.AttachmentServerRemoveUploaderChan <- connData.Uid
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

		socketServer.CallsPending.mutex.Lock()
		if callPending, ok := socketServer.CallsPending.data[connData.Uid]; ok {
			socketServer.SendDataToUser <- UserDataMessage{
				Uid:  callPending,
				Type: "CALL_USER_RESPONSE",
				Data: socketmodels.CallResponse{
					Caller: connData.Uid.Hex(),
					Called: callPending.Hex(),
					Accept: false,
				},
			}
			delete(socketServer.CallsActive.data, connData.Uid)
		}
		for caller, called := range socketServer.CallsPending.data {
			if called == connData.Uid {
				socketServer.SendDataToUser <- UserDataMessage{
					Uid:  caller,
					Type: "CALL_USER_RESPONSE",
					Data: socketmodels.CallResponse{
						Caller: caller.Hex(),
						Called: connData.Uid.Hex(),
						Accept: false,
					},
				}
				delete(socketServer.CallsPending.data, caller)
			}
		}
		socketServer.CallsPending.mutex.Unlock()

		socketServer.CallsActive.mutex.Lock()
		if called, ok := socketServer.CallsActive.data[connData.Uid]; ok {
			socketServer.SendDataToUser <- UserDataMessage{
				Uid:  called,
				Type: "CALL_LEFT",
				Data: socketmodels.CallLeft{},
			}
			delete(socketServer.CallsActive.data, connData.Uid)
		} else {
			for caller, called := range socketServer.CallsActive.data {
				if called == connData.Uid {
					socketServer.SendDataToUser <- UserDataMessage{
						Type: "CALL_LEFT",
						Uid:  caller,
						Data: socketmodels.CallLeft{},
					}
					delete(socketServer.CallsActive.data, caller)
					break
				}
			}
		}
		socketServer.CallsActive.mutex.Unlock()

		socketServer.Connections.mutex.Unlock()
		socketServer.Subscriptions.mutex.Unlock()
	}
}

func messageQueueLoop(socketServer *SocketServer, colls *db.Collections) {
	for {
		defer func() {
			r := recover()
			if r != nil {
				log.Println("Recovered from panic in queued socket messages :", r)
			}
			go messageQueueLoop(socketServer, colls)
		}()
		data := <-socketServer.MessageSendQueue
		data.Conn.WriteMessage(websocket.TextMessage, data.Data)
	}
}

func subscriptionConnectionRegistrationLoop(socketServer *SocketServer, colls *db.Collections) {
	for {
		defer func() {
			r := recover()
			if r != nil {
				log.Println("Recovered from panic in subscription registration :", r)
			}
			go subscriptionConnectionRegistrationLoop(socketServer, colls)
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
}

func subscriptionDisconnectRegistrationLoop(socketServer *SocketServer, colls *db.Collections) {
	for {
		defer func() {
			r := recover()
			if r != nil {
				log.Println("Recovered from panic in subscription disconnect registration :", r)
			}
			go subscriptionDisconnectRegistrationLoop(socketServer, colls)
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
}

func sendSubscriptionDataLoop(socketServer *SocketServer, colls *db.Collections) {
	for {
		defer func() {
			r := recover()
			if r != nil {
				log.Println("Recovered from panic in subscription data channel:", r)
			}
			go sendSubscriptionDataLoop(socketServer, colls)
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
}

func sendSubscriptionDataExclusiveLoop(socketServer *SocketServer, colls *db.Collections) {
	for {
		defer func() {
			r := recover()
			if r != nil {
				log.Println("Recovered from panic in exclusive subscription data channel:", r)
			}
			go sendSubscriptionDataExclusiveLoop(socketServer, colls)
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
}

func sendToMultipleSubscriptionsLoop(socketServer *SocketServer, colls *db.Collections) {
	for {
		defer func() {
			r := recover()
			if r != nil {
				log.Println("Recovered from panic in subscription data channel:", r)
			}
			go sendToMultipleSubscriptionsLoop(socketServer, colls)
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
}

func sendToMultipleSubscriptionsExclusiveLoop(socketServer *SocketServer, colls *db.Collections) {
	for {
		defer func() {
			r := recover()
			if r != nil {
				log.Println("Recovered from panic in exclusive subscription data channel:", r)
			}
			go sendToMultipleSubscriptionsExclusiveLoop(socketServer, colls)
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
}

func getSubscriptionUidsLoop(socketServer *SocketServer, colls *db.Collections) {
	for {
		defer func() {
			r := recover()
			if r != nil {
				log.Println("Recovered from panic getting uids from subscription channel:", r)
			}
			go getSubscriptionUidsLoop(socketServer, colls)
		}()
		subsData := <-socketServer.GetSubscriptionUids
		socketServer.Subscriptions.mutex.Lock()
		uids := make(map[primitive.ObjectID]struct{})
		for _, oi := range socketServer.Subscriptions.data[subsData.Name] {
			uids[oi] = struct{}{}
		}
		subsData.RecvChan <- uids
		socketServer.Subscriptions.mutex.Unlock()
	}
}

func sendDataToUserLoop(socketServer *SocketServer, colls *db.Collections) {
	for {
		defer func() {
			r := recover()
			if r != nil {
				log.Println("Recovered from panic in send data to user channel:", r)
			}
			go sendDataToUserLoop(socketServer, colls)
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
}

func sendDataToUsersLoop(socketServer *SocketServer, colls *db.Collections) {
	for {
		defer func() {
			r := recover()
			if r != nil {
				log.Println("Recovered from panic in send data to users channel:", r)
			}
			go sendDataToUsersLoop(socketServer, colls)
		}()
		data := <-socketServer.SendDataToUsers
		socketServer.Connections.mutex.Lock()
		m := make(map[string]interface{})
		outBytesNoTypeKey, err := json.Marshal(data.Data)
		json.Unmarshal(outBytesNoTypeKey, &m)
		m["TYPE"] = data.Type
		outBytes, err := json.Marshal(m)
		if err != nil {
			log.Println("Error marshaling data to be sent to user :", err)
			continue
		}
		for conn, uid := range socketServer.Connections.data {
			_, ok := data.Uids[uid]
			if ok {
				socketServer.MessageSendQueue <- QueuedMessage{
					Conn: conn,
					Data: outBytes,
				}
			}
		}
		socketServer.Connections.mutex.Unlock()
	}
}

func removeUserFromSubscriptionLoop(socketServer *SocketServer, colls *db.Collections) {
	for {
		defer func() {
			r := recover()
			if r != nil {
				log.Println("Recovered from panic in remove user from subscription channel:", r)
			}
			go removeUserFromSubscriptionLoop(socketServer, colls)
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
}

func destroySubscriptionLoop(socketServer *SocketServer, colls *db.Collections) {
	for {
		defer func() {
			r := recover()
			if r != nil {
				log.Println("Recovered from panic in destroy subscription channel:", r)
			}
			go destroySubscriptionLoop(socketServer, colls)
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
}

func callPendingChanLoop(socketServer *SocketServer, colls *db.Collections) {
	for {
		defer func() {
			r := recover()
			if r != nil {
				log.Println("Recovered from panic in pending call channel:", r)
			}
			go callPendingChanLoop(socketServer, colls)
		}()
		data := <-socketServer.CallsPendingChan
		socketServer.CallsPending.mutex.Lock()
		if called, ok := socketServer.CallsPending.data[data.Caller]; ok {
			if called != data.Called {
				// pending call switching to different user. cancel previous pending call.
				Uids := make(map[primitive.ObjectID]struct{})
				Uids[called] = struct{}{}
				Uids[data.Caller] = struct{}{}
				socketServer.SendDataToUsers <- UsersDataMessage{
					Uids: Uids,
					Type: "CALL_USER_RESPONSE",
					Data: socketmodels.CallResponse{
						Called: data.Called.Hex(),
						Caller: data.Caller.Hex(),
						Accept: false,
					},
				}
				socketServer.CallsPending.data[data.Caller] = data.Called
			}
		} else {
			socketServer.CallsPending.data[data.Caller] = data.Called
		}
		Uids := make(map[primitive.ObjectID]struct{})
		Uids[data.Called] = struct{}{}
		Uids[data.Caller] = struct{}{}
		socketServer.SendDataToUsers <- UsersDataMessage{
			Uids: Uids,
			Type: "CALL_USER_ACKNOWLEDGE",
			Data: socketmodels.CallAcknowledge{
				Caller: data.Caller.Hex(),
				Called: data.Called.Hex(),
			},
		}
		socketServer.CallsPending.mutex.Unlock()
	}
}

func callResponseChanLoop(socketServer *SocketServer, colls *db.Collections) {
	for {
		defer func() {
			r := recover()
			if r != nil {
				log.Println("Recovered from panic in call response channel:", r)
			}
			go callResponseChanLoop(socketServer, colls)
		}()
		data := <-socketServer.ResponseToCallChan
		socketServer.CallsPending.mutex.Lock()
		socketServer.CallsActive.mutex.Lock()
		delete(socketServer.CallsPending.data, data.Caller)

		if data.Accept {
			// Close any call that either user is currently in.
			// Clients can only be in a single call.
			// Confusing variable names here.
			closedCallerCall := false
			closedCalledCall := false
			if callerCalled, ok := socketServer.CallsActive.data[data.Caller]; ok {
				closedCallerCall = true
				Uids := make(map[primitive.ObjectID]struct{})
				Uids[data.Caller] = struct{}{}
				Uids[callerCalled] = struct{}{}
				socketServer.SendDataToUsers <- UsersDataMessage{
					Type: "CALL_LEFT",
					Data: socketmodels.CallLeft{},
					Uids: Uids,
				}
				delete(socketServer.CallsActive.data, data.Caller)
			}
			if calledCalled, ok := socketServer.CallsActive.data[data.Called]; ok {
				closedCalledCall = true
				Uids := make(map[primitive.ObjectID]struct{})
				Uids[data.Called] = struct{}{}
				Uids[calledCalled] = struct{}{}
				socketServer.SendDataToUsers <- UsersDataMessage{
					Type: "CALL_LEFT",
					Data: socketmodels.CallLeft{},
					Uids: Uids,
				}
				delete(socketServer.CallsActive.data, data.Called)
			}
			// make sure that the caller is not in a call. If they are exit the call they are already in
			if !closedCallerCall {
				for caller, called := range socketServer.CallsActive.data {
					if data.Caller == called {
						Uids := make(map[primitive.ObjectID]struct{})
						Uids[caller] = struct{}{}
						Uids[called] = struct{}{}
						socketServer.SendDataToUsers <- UsersDataMessage{
							Type: "CALL_LEFT",
							Data: socketmodels.CallLeft{},
							Uids: Uids,
						}
						delete(socketServer.CallsActive.data, caller)
						break
					}
				}
			}
			// make sure that the called user is not in a call. If they are exit the call they are already in
			if !closedCalledCall {
				for caller, called := range socketServer.CallsActive.data {
					if data.Called == called {
						Uids := make(map[primitive.ObjectID]struct{})
						Uids[caller] = struct{}{}
						Uids[called] = struct{}{}
						socketServer.SendDataToUsers <- UsersDataMessage{
							Type: "CALL_LEFT",
							Data: socketmodels.CallLeft{},
							Uids: Uids,
						}
						delete(socketServer.CallsActive.data, caller)
						break
					}
				}
			}

			// Any active calls that either user in have now been closed. Proceed.
			socketServer.CallsActive.data[data.Caller] = data.Called
		}

		// Send the response to both clients
		Uids := make(map[primitive.ObjectID]struct{})
		Uids[data.Called] = struct{}{}
		Uids[data.Caller] = struct{}{}
		socketServer.SendDataToUsers <- UsersDataMessage{
			Uids: Uids,
			Type: "CALL_USER_RESPONSE",
			Data: socketmodels.CallResponse{
				Caller: data.Caller.Hex(),
				Called: data.Called.Hex(),
				Accept: data.Accept,
			},
		}

		socketServer.CallsPending.mutex.Unlock()
		socketServer.CallsActive.mutex.Unlock()
	}
}
