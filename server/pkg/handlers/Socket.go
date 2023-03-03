package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/web-stuff-98/electron-social-chat/pkg/db"
	"github.com/web-stuff-98/electron-social-chat/pkg/helpers"
	"github.com/web-stuff-98/electron-social-chat/pkg/socketserver"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
}

/*
	Socket event handling.

	Voting and commenting are done in the API handlers, I could have put that in here but I didn't

	Todo:
	 - sendErrorMessageThroughSocket with http status code and message
*/

func reader(conn *websocket.Conn, socketServer *socketserver.SocketServer, uid *primitive.ObjectID, colls *db.Collections) {
	for {
		defer func() {
			r := recover()
			if r != nil {
				log.Println("Recovered from panic in WS reader loop : ", r)
			}
		}()

		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		var data map[string]interface{}
		json.Unmarshal(p, &data)

		eventType, eventTypeOk := data["event_type"]

		if eventTypeOk {
			err := HandleSocketEvent(eventType.(string), p, conn, *uid, socketServer, colls)
			if err != nil {
				sendErrorMessageThroughSocket(conn, err)
			}
		} else {
			// eventType was not received. Send error.
			sendErrorMessageThroughSocket(conn, err)
		}
	}
}

func sendErrorMessageThroughSocket(conn *websocket.Conn, e error) {
	err := conn.WriteJSON(map[string]string{
		"TYPE": "RESPONSE_MESSAGE",
		"DATA": `{"msg":"` + e.Error() + `","err":true}`,
	})
	if err != nil {
		log.Println(err)
	}
}

func (h handler) WebSocketEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	user, err := helpers.GetUserFromRequest(r, context.Background(), *h.Collections, h.RedisClient)
	uid := primitive.NilObjectID
	if user != nil {
		uid = user.ID
	}
	h.SocketServer.RegisterConn <- socketserver.ConnectionInfo{
		Conn:   ws,
		Uid:    uid,
		Online: true,
	}
	defer func() {
		h.SocketServer.UnregisterConn <- socketserver.ConnectionInfo{
			Conn:   ws,
			Uid:    uid,
			Online: false,
		}
	}()
	log.Println("Connection")
	reader(ws, h.SocketServer, &uid, h.Collections)
}
