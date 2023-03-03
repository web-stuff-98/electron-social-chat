package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/web-stuff-98/electron-social-chat/pkg/db"
	"github.com/web-stuff-98/electron-social-chat/pkg/db/models"
	"github.com/web-stuff-98/electron-social-chat/pkg/socketmodels"
	"github.com/web-stuff-98/electron-social-chat/pkg/socketserver"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
	case "ROOM_OPEN_CHANNEL":
		err := openRoomChannel(data, conn, uid, ss, colls)
		return err
	case "ROOM_EXIT_CHANNEL":
		err := exitRoomChannel(data, conn, uid, ss, colls)
		return err
	case "ROOM_MESSAGE":
		err := roomMessage(data, conn, uid, ss, colls)
		return err
	case "ROOM_MESSAGE_UPDATE":
		err := roomMessageUpdate(data, conn, uid, ss, colls)
		return err
	case "ROOM_MESSAGE_DELETE":
		err := roomMessageDelete(data, conn, uid, ss, colls)
		return err
	}

	return fmt.Errorf("Unrecognized event type")
}

func openSubscription(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	var data socketmodels.OpenCloseSubscription
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	if strings.HasPrefix(data.Name, "channel:") {
		return fmt.Errorf("Unauthorized")
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
		if strings.HasPrefix(name, "channel:") {
			continue
		}
		ss.RegisterSubscriptionConn <- socketserver.SubscriptionConnectionInfo{
			Name: name,
			Uid:  uid,
			Conn: conn,
		}
	}
	return nil
}

func openRoomChannel(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	var data socketmodels.RoomOpenExitChannel
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	channelId, err := primitive.ObjectIDFromHex(data.Channel)
	if err != nil {
		return err
	}

	channel := &models.RoomChannel{}
	if err := colls.RoomChannelCollection.FindOne(context.Background(), bson.M{"_id": channelId}).Decode(&channel); err != nil {
		return err
	}
	roomExternalData := &models.RoomExternalData{}
	if err := colls.RoomExternalDataCollection.FindOne(context.Background(), bson.M{"_id": channel.RoomID}).Decode(&roomExternalData); err != nil {
		return err
	}

	for _, oi := range roomExternalData.Banned {
		if oi == uid {
			return fmt.Errorf("Banned")
		}
	}
	if roomExternalData.Private {
		member := false
		for _, oi := range roomExternalData.Members {
			if oi == uid {
				member = true
				break
			}
		}
		if !member {
			return fmt.Errorf("Not a member")
		}
	}

	ss.RegisterSubscriptionConn <- socketserver.SubscriptionConnectionInfo{
		Name: "channel:" + channelId.Hex(),
		Uid:  uid,
		Conn: conn,
	}

	return nil
}

func exitRoomChannel(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	var data socketmodels.RoomOpenExitChannel
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	channelId, err := primitive.ObjectIDFromHex(data.Channel)
	if err != nil {
		return err
	}

	channel := &models.RoomChannel{}
	if err := colls.RoomChannelCollection.FindOne(context.Background(), bson.M{"_id": channelId}).Decode(&channel); err != nil {
		return err
	}
	roomExternalData := &models.RoomExternalData{}
	if err := colls.RoomExternalDataCollection.FindOne(context.Background(), bson.M{"_id": channel.RoomID}).Decode(&roomExternalData); err != nil {
		return err
	}

	ss.UnregisterSubscriptionConn <- socketserver.SubscriptionConnectionInfo{
		Name: "channel:" + channelId.Hex(),
		Uid:  uid,
		Conn: conn,
	}

	return nil
}

func roomMessage(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	var data socketmodels.RoomMessage
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	if strings.TrimSpace(data.Content) == "" {
		return fmt.Errorf("You cannot submit an empty message")
	}
	if len(data.Content) > 200 {
		return fmt.Errorf("Max 200 characters")
	}

	channelId, err := primitive.ObjectIDFromHex(data.Channel)
	if err != nil {
		return err
	}

	channel := &models.RoomChannel{}
	if err := colls.RoomChannelCollection.FindOne(context.Background(), bson.M{"_id": channelId}).Decode(&channel); err != nil {
		return err
	}
	roomExternalData := &models.RoomExternalData{}
	if err := colls.RoomExternalDataCollection.FindOne(context.Background(), bson.M{"_id": channel.RoomID}).Decode(&roomExternalData); err != nil {
		return err
	}

	for _, oi := range roomExternalData.Banned {
		if oi == uid {
			return fmt.Errorf("Banned")
		}
	}
	if roomExternalData.Private {
		member := false
		for _, oi := range roomExternalData.Members {
			if oi == uid {
				member = true
				break
			}
		}
		if !member {
			return fmt.Errorf("Not a member")
		}
	}

	msgId := primitive.NewObjectID()

	if _, err := colls.RoomChannelMessagesCollection.UpdateByID(context.Background(), bson.M{"_id": channel.ID}, bson.M{
		"messages": bson.M{
			"$push": models.RoomChannelMessage{
				ID:        msgId,
				Content:   data.Content,
				CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
				UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
				Author:    uid,
			},
		},
	}); err != nil {
		return err
	}

	outBytes, err := json.Marshal(socketmodels.OutRoomMessage{
		Type:    "OUT_ROOM_MESSAGE",
		Content: data.Content,
		ID:      msgId.Hex(),
		Author:  uid.Hex(),
	})

	ss.SendDataToSubscription <- socketserver.SubscriptionDataMessage{
		Name: "channel:" + channelId.Hex(),
		Data: outBytes,
	}

	return nil
}

func roomMessageUpdate(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	var data socketmodels.RoomMessageUpdate
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	if strings.TrimSpace(data.Content) == "" {
		return fmt.Errorf("You cannot submit an empty message")
	}
	if len(data.Content) > 200 {
		return fmt.Errorf("Max 200 characters")
	}

	channelId, err := primitive.ObjectIDFromHex(data.Channel)
	if err != nil {
		return err
	}
	msgId, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return err
	}

	channel := &models.RoomChannel{}
	if err := colls.RoomChannelCollection.FindOne(context.Background(), bson.M{"_id": channelId}).Decode(&channel); err != nil {
		return err
	}
	roomExternalData := &models.RoomExternalData{}
	if err := colls.RoomExternalDataCollection.FindOne(context.Background(), bson.M{"_id": channel.RoomID}).Decode(&roomExternalData); err != nil {
		return err
	}

	for _, oi := range roomExternalData.Banned {
		if oi == uid {
			return fmt.Errorf("Banned")
		}
	}
	if roomExternalData.Private {
		member := false
		for _, oi := range roomExternalData.Members {
			if oi == uid {
				member = true
				break
			}
		}
		if !member {
			return fmt.Errorf("Not a member")
		}
	}

	if res, err := colls.RoomChannelMessagesCollection.UpdateByID(context.Background(), bson.M{
		"_id":             channelId,
		"messages._id":    msgId,
		"messages.author": uid,
	}, bson.M{
		"$set": bson.M{
			"messages.$.content": data.Content,
		},
	}); err != nil {
		return err
	} else if res.MatchedCount == 0 {
		return fmt.Errorf("Update failed")
	}

	outBytes, err := json.Marshal(socketmodels.OutRoomMessageUpdate{
		Type:    "OUT_ROOM_MESSAGE_UPDATE",
		Content: data.Content,
		ID:      msgId.Hex(),
	})

	ss.SendDataToSubscription <- socketserver.SubscriptionDataMessage{
		Name: "channel:" + channelId.Hex(),
		Data: outBytes,
	}

	return nil
}

func roomMessageDelete(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	var data socketmodels.RoomMessageDelete
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	channelId, err := primitive.ObjectIDFromHex(data.Channel)
	if err != nil {
		return err
	}
	msgId, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return err
	}

	channel := &models.RoomChannel{}
	if err := colls.RoomChannelCollection.FindOne(context.Background(), bson.M{"_id": channelId}).Decode(&channel); err != nil {
		return err
	}
	roomExternalData := &models.RoomExternalData{}
	if err := colls.RoomExternalDataCollection.FindOne(context.Background(), bson.M{"_id": channel.RoomID}).Decode(&roomExternalData); err != nil {
		return err
	}

	for _, oi := range roomExternalData.Banned {
		if oi == uid {
			return fmt.Errorf("Banned")
		}
	}
	if roomExternalData.Private {
		member := false
		for _, oi := range roomExternalData.Members {
			if oi == uid {
				member = true
				break
			}
		}
		if !member {
			return fmt.Errorf("Not a member")
		}
	}

	if res, err := colls.RoomChannelMessagesCollection.UpdateByID(context.Background(), bson.M{
		"_id": channelId,
	}, bson.M{
		"$pull": bson.M{
			"messages": bson.M{
				"_id":    msgId,
				"author": uid,
			},
		},
	}); err != nil {
		return err
	} else if res.MatchedCount == 0 {
		return fmt.Errorf("Delete failed")
	}

	outBytes, err := json.Marshal(socketmodels.OutRoomMessageDelete{
		Type: "OUT_ROOM_MESSAGE_DELETE",
		ID:   msgId.Hex(),
	})

	ss.SendDataToSubscription <- socketserver.SubscriptionDataMessage{
		Name: "channel:" + channelId.Hex(),
		Data: outBytes,
	}

	return nil
}
