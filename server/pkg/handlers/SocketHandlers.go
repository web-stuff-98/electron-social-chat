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
	"go.mongodb.org/mongo-driver/mongo/options"
)

func HandleSocketEvent(eventType string, data []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	switch eventType {
	case "WATCH_USER":
		err := watchUser(data, conn, uid, ss, colls)
		return err
	case "STOP_WATCHING_USER":
		err := stopWatchingUser(data, conn, uid, ss, colls)
		return err
	case "WATCH_ROOM":
		err := watchRoom(data, conn, uid, ss, colls)
		return err
	case "STOP_WATCHING_ROOM":
		err := stopWatchingRoom(data, conn, uid, ss, colls)
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
	case "DIRECT_MESSAGE":
		err := directMessage(data, conn, uid, ss, colls)
		return err
	case "DIRECT_MESSAGE_UPDATE":
		err := directMessageUpdate(data, conn, uid, ss, colls)
		return err
	case "DIRECT_MESSAGE_DELETE":
		err := directMessageDelete(data, conn, uid, ss, colls)
		return err
	}
	return fmt.Errorf("Unrecognized event type")
}

func watchUser(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	var data socketmodels.WatchStopWatching
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	id, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return err
	}

	ss.RegisterSubscriptionConn <- socketserver.SubscriptionConnectionInfo{
		Name: "user=" + id.Hex(),
		Uid:  uid,
		Conn: conn,
	}

	return nil
}

func stopWatchingUser(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	var data socketmodels.WatchStopWatching
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	id, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return err
	}

	ss.UnregisterSubscriptionConn <- socketserver.SubscriptionConnectionInfo{
		Name: "user=" + id.Hex(),
		Uid:  uid,
		Conn: conn,
	}

	return nil
}

func watchRoom(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	var data socketmodels.WatchStopWatching
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	id, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return err
	}

	ss.RegisterSubscriptionConn <- socketserver.SubscriptionConnectionInfo{
		Name: "room-display-data=" + id.Hex(),
		Uid:  uid,
		Conn: conn,
	}

	return nil
}

func stopWatchingRoom(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	var data socketmodels.WatchStopWatching
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	id, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return err
	}

	ss.UnregisterSubscriptionConn <- socketserver.SubscriptionConnectionInfo{
		Name: "room-display-data=" + id.Hex(),
		Uid:  uid,
		Conn: conn,
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
	if len(data.Content) > 300 {
		return fmt.Errorf("Max 300 characters")
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

	if _, err := colls.RoomChannelMessagesCollection.UpdateByID(context.Background(), channel.ID, bson.M{
		"$push": bson.M{
			"messages": models.RoomChannelMessage{
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
	if len(data.Content) > 300 {
		return fmt.Errorf("Max 300 characters")
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

	if res, err := colls.RoomChannelMessagesCollection.UpdateOne(context.Background(), bson.M{
		"_id":             channelId,
		"messages._id":    msgId,
		"messages.author": uid,
	}, bson.M{
		"$set": bson.M{
			"messages.$.content":    data.Content,
			"messages.$.updated_at": primitive.NewDateTimeFromTime(time.Now()),
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

	if err != nil {
		return err
	}

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

	if res, err := colls.RoomChannelMessagesCollection.UpdateByID(context.Background(), channelId, bson.M{
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

func directMessage(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	var data socketmodels.DirectMessage
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	recipientId, err := primitive.ObjectIDFromHex(data.Recipient)
	if err != nil {
		return err
	}

	recipientMessagingData := &models.UserMessagingData{}
	if err := colls.UserMessagingDataCollection.FindOne(context.Background(), bson.M{"_id": recipientId}).Decode(&recipientMessagingData); err != nil {
		return err
	} else {
		for _, oi := range recipientMessagingData.Blocked {
			if oi == uid {
				return fmt.Errorf("This user has blocked your account")
			}
		}
	}

	msgId := primitive.NewObjectID()

	if _, err := colls.UserMessagingDataCollection.UpdateByID(context.Background(), recipientId, bson.M{
		"$push": bson.M{
			"messages": models.DirectMessage{
				ID:        msgId,
				CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
				UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
				Author:    uid,
			},
		},
		"$addToSet": bson.M{
			"messages_received_from": uid,
		},
	}); err != nil {
		return err
	}

	if _, err := colls.UserMessagingDataCollection.UpdateByID(context.Background(), uid, bson.M{
		"$addToSet": bson.M{
			"messages_sent_to": recipientId,
		},
	}); err != nil {
		return err
	}

	msg := &socketmodels.OutDirectMessage{
		ID:        msgId.Hex(),
		Content:   data.Content,
		Author:    uid.Hex(),
		Recipient: recipientId.Hex(),
	}
	ss.SendDataToUser <- socketserver.UserDataMessage{
		Uid:  uid,
		Type: "OUT_DIRECT_MESSAGE",
		Data: msg,
	}
	ss.SendDataToUser <- socketserver.UserDataMessage{
		Uid:  recipientId,
		Type: "OUT_DIRECT_MESSAGE",
		Data: msg,
	}

	return nil
}

func directMessageUpdate(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	var data socketmodels.DirectMessageUpdate
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	recipientId, err := primitive.ObjectIDFromHex(data.Recipient)
	if err != nil {
		return err
	}

	msgId, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return err
	}

	if res, err := colls.UserMessagingDataCollection.UpdateOne(context.Background(), bson.M{
		"_id":             recipientId,
		"messages._id":    msgId,
		"messages.author": uid,
	}, bson.M{
		"$set": bson.M{
			"messages.$.content":    data.Content,
			"messages.$.updated_at": primitive.NewDateTimeFromTime(time.Now()),
		},
	}); err != nil {
		return err
	} else if res.MatchedCount == 0 {
		return fmt.Errorf("Update failed")
	}

	msg := &socketmodels.OutDirectMessageUpdate{
		ID:        msgId.Hex(),
		Content:   data.Content,
		Author:    uid.Hex(),
		Recipient: recipientId.Hex(),
	}
	ss.SendDataToUser <- socketserver.UserDataMessage{
		Uid:  uid,
		Type: "OUT_DIRECT_MESSAGE_UPDATE",
		Data: msg,
	}
	ss.SendDataToUser <- socketserver.UserDataMessage{
		Uid:  recipientId,
		Type: "OUT_DIRECT_MESSAGE_UPDATE",
		Data: msg,
	}

	return nil
}

func directMessageDelete(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	var data socketmodels.DirectMessageDelete
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	recipientId, err := primitive.ObjectIDFromHex(data.Recipient)
	if err != nil {
		return err
	}

	msgId, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return err
	}

	recipientMessagingData := &models.UserMessagingData{}
	if err := colls.UserMessagingDataCollection.FindOneAndUpdate(context.Background(), bson.M{
		"_id": recipientId,
	}, bson.M{
		"$pull": bson.M{
			"messages": bson.M{
				"_id":    msgId,
				"author": uid,
			},
		},
	}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&recipientMessagingData); err != nil {
		return err
	} else {
		wasOnlyMessage := true
		for _, dm := range recipientMessagingData.Messages {
			if dm.Author == uid {
				wasOnlyMessage = false
				break
			}
		}
		if wasOnlyMessage {
			if _, err := colls.UserMessagingDataCollection.UpdateByID(context.Background(), recipientId, bson.M{
				"$pull": bson.M{
					"messages_received_from": uid,
				},
			}); err != nil {
				return err
			}
			if _, err := colls.UserMessagingDataCollection.UpdateByID(context.Background(), uid, bson.M{
				"$pull": bson.M{
					"messages_sent_to": uid,
				},
			}); err != nil {
				return err
			}
		}
	}

	msg := &socketmodels.OutDirectMessageDelete{
		ID:        msgId.Hex(),
		Author:    uid.Hex(),
		Recipient: recipientId.Hex(),
	}
	ss.SendDataToUser <- socketserver.UserDataMessage{
		Uid:  uid,
		Type: "OUT_DIRECT_MESSAGE_DELETE",
		Data: msg,
	}
	ss.SendDataToUser <- socketserver.UserDataMessage{
		Uid:  recipientId,
		Type: "OUT_DIRECT_MESSAGE_DELETE",
		Data: msg,
	}

	return nil
}
