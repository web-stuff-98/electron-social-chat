package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/web-stuff-98/electron-social-chat/pkg/db"
	"github.com/web-stuff-98/electron-social-chat/pkg/socketmodels"
	"github.com/web-stuff-98/electron-social-chat/pkg/socketserver"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
	Moved from Socket.go, because the code was using tonnes of if/else statements for error
	handling and was getting messy looking.
*/

func HandleSocketEvent(eventType string, data []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	switch eventType {
	case "OPEN_SUBSCRIPTION":
		err := openSubscription(data, conn, uid, ss, colls)
		return err
	case "CLOSE_SUBSCRIPTION":
		err := closeSubscription(data, conn, uid, ss, colls)
		return err
	case "OPEN_SUBSCRIPTIONS":
		err := openSubscriptions(data, conn, uid, ss, colls)
		return err
	}
	return fmt.Errorf("Unrecognized event type")
}

func openSubscription(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	var data socketmodels.OpenCloseSubscription
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	ss.RegisterSubscriptionConn <- socketserver.SubscriptionConnectionInfo{
		Name: data.Name,
		Uid:  uid,
		Conn: conn,
	}
	return nil
}

func closeSubscription(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	var data socketmodels.OpenCloseSubscription
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	ss.UnregisterSubscriptionConn <- socketserver.SubscriptionConnectionInfo{
		Name: data.Name,
		Uid:  uid,
		Conn: conn,
	}
	return nil
}

func openSubscriptions(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	var data socketmodels.OpenCloseSubscriptions
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	for _, name := range data.Names {
		ss.RegisterSubscriptionConn <- socketserver.SubscriptionConnectionInfo{
			Name: name,
			Uid:  uid,
			Conn: conn,
		}
	}
	return nil
}
