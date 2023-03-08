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
	case "FRIEND_REQUEST":
		err := friendRequest(data, conn, uid, ss, colls)
		return err
	case "FRIEND_REQUEST_RESPONSE":
		err := friendRequestResponse(data, conn, uid, ss, colls)
		return err
	case "BLOCK_USER":
		err := blockUser(data, conn, uid, ss, colls)
		return err
	case "UNBLOCK_USER":
		err := unblockUser(data, conn, uid, ss, colls)
		return err
	case "ROOM_INVITATION":
		err := inviteToRoom(data, conn, uid, ss, colls)
		return err
	case "ROOM_INVITATION_RESPONSE":
		err := invitationToRoomResponse(data, conn, uid, ss, colls)
		return err
	case "ROOM_INVITATION_DELETE":
		err := deleteInvitationToRoom(data, conn, uid, ss, colls)
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
	room := &models.Room{}
	if err := colls.RoomCollection.FindOne(context.Background(), bson.M{"_id": channel.RoomID}).Decode(&room); err != nil {
		return err
	}

	if uid != room.Author {
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
	room := &models.Room{}
	if err := colls.RoomCollection.FindOne(context.Background(), bson.M{"_id": channel.RoomID}).Decode(&room); err != nil {
		return err
	}

	if room.Author != uid {
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
	room := &models.Room{}
	if err := colls.RoomChannelCollection.FindOne(context.Background(), bson.M{"_id": channel.RoomID}).Decode(&room); err != nil {
		return err
	}

	if room.Author != uid {
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
	room := &models.Room{}
	if err := colls.RoomCollection.FindOne(context.Background(), bson.M{"_id": channel.RoomID}).Decode(&room); err != nil {
		return err
	}

	if room.Author != uid {
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
				Content:   data.Content,
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
		wasOnlyMessage := checkAnythingReceivedFrom(*recipientMessagingData, uid)
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
	Uids := make(map[primitive.ObjectID]struct{})
	Uids[uid] = struct{}{}
	Uids[recipientId] = struct{}{}
	ss.SendDataToUsers <- socketserver.UsersDataMessage{
		Uids: Uids,
		Type: "OUT_DIRECT_MESSAGE_DELETE",
		Data: msg,
	}

	return nil
}

func inviteToRoom(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	var data socketmodels.InviteToRoom
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	recipientId, err := primitive.ObjectIDFromHex(data.Recipient)
	if err != nil {
		return err
	}
	roomId, err := primitive.ObjectIDFromHex(data.RoomID)
	if err != nil {
		return err
	}

	room := &models.Room{}
	if err := colls.RoomCollection.FindOne(context.Background(), bson.M{"_id": roomId}).Decode(&room); err != nil {
		return err
	}

	roomExternalData := &models.RoomExternalData{}
	if err := colls.RoomExternalDataCollection.FindOne(context.Background(), bson.M{"_id": roomId}).Decode(&roomExternalData); err != nil {
		return err
	}

	for _, oi := range roomExternalData.Banned {
		if oi == recipientId {
			return fmt.Errorf("You have banned this user. You must unban them to send an invite")
		}
	}
	for _, oi := range roomExternalData.Members {
		if oi == recipientId {
			return fmt.Errorf("This user is already a member of the room")
		}
	}

	if room.Author != uid {
		return fmt.Errorf("Unauthorized")
	}

	messagingData := &models.UserMessagingData{}
	if err := colls.UserMessagingDataCollection.FindOne(context.Background(), bson.M{"_id": recipientId}).Decode(&messagingData); err != nil {
		return err
	}
	for _, oi := range messagingData.Blocked {
		if oi == recipientId {
			return fmt.Errorf("You have blocked this users account. You must unblock the user before sending them an invite")
		}
	}

	recipientMessagingData := &models.UserMessagingData{}
	if err := colls.UserMessagingDataCollection.FindOne(context.Background(), bson.M{"_id": recipientId}).Decode(&recipientMessagingData); err != nil {
		return err
	} else {
		for _, oi := range recipientMessagingData.Blocked {
			if oi == uid {
				return fmt.Errorf("This user has blocked your account, you cannot send them invitations")
			}
		}
	}

	invitationId := primitive.NewObjectID()

	if _, err := colls.UserMessagingDataCollection.UpdateByID(context.Background(), recipientId, bson.M{
		"$push": bson.M{
			"invitations": models.Invitation{
				ID:        invitationId,
				CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
				Author:    uid,
				RoomID:    roomId,
				Accepted:  false,
				Declined:  false,
			},
		},
		"$addToSet": bson.M{
			// invitations count as messages, even though they are stored in a seperate array
			"messages_received_from": uid,
		},
	}); err != nil {
		return err
	}

	if _, err := colls.UserMessagingDataCollection.UpdateByID(context.Background(), uid, bson.M{
		"$addToSet": bson.M{
			// invitations count as messages, even though they are stored in a seperate array
			"messages_sent_to": recipientId,
		},
	}); err != nil {
		return err
	}

	msg := &socketmodels.OutInvite{
		ID:        invitationId.Hex(),
		Author:    uid.Hex(),
		Recipient: recipientId.Hex(),
		RoomID:    roomId.Hex(),
	}
	Uids := make(map[primitive.ObjectID]struct{})
	Uids[uid] = struct{}{}
	Uids[recipientId] = struct{}{}
	ss.SendDataToUsers <- socketserver.UsersDataMessage{
		Uids: Uids,
		Type: "OUT_DIRECT_MESSAGE",
		Data: msg,
	}

	return nil
}

func deleteInvitationToRoom(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	var data socketmodels.RoomInvitationDelete
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	invitationId, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return err
	}

	messagingData := &models.UserMessagingData{}
	if err := colls.UserMessagingDataCollection.FindOne(context.Background(), bson.M{"_id": uid}).Decode(&messagingData); err != nil {
		return err
	}

	invitationIndex := -1
	for i, inv := range messagingData.Invitations {
		if inv.ID == invitationId {
			invitationIndex = i
			break
		}
	}
	if invitationIndex == -1 {
		return fmt.Errorf("Invitation not found")
	}

	if messagingData.Invitations[invitationIndex].Author != uid {
		return fmt.Errorf("Unauthorized")
	}

	if err := colls.UserMessagingDataCollection.FindOneAndUpdate(context.Background(), bson.M{"_id": uid}, bson.M{
		"$pull": bson.M{
			"invitations": bson.M{
				"_id":    invitationId,
				"author": messagingData.Invitations[invitationIndex].Author,
			},
		},
	}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&messagingData); err != nil {
		return err
	} else {
		if !checkAnythingReceivedFrom(*messagingData, messagingData.Invitations[invitationIndex].Author) {
			if _, err := colls.UserMessagingDataCollection.UpdateByID(context.Background(), uid, bson.M{
				"$pull": bson.M{
					"messages_received_from": messagingData.Invitations[invitationIndex].Author,
				},
			}); err != nil {
				return err
			}
		}
	}

	msg := &socketmodels.OutRoomInvitationDelete{
		ID:     invitationId.Hex(),
		Author: uid.Hex(),
	}
	Uids := make(map[primitive.ObjectID]struct{})
	Uids[uid] = struct{}{}
	Uids[messagingData.Invitations[invitationIndex].Author] = struct{}{}
	ss.SendDataToUsers <- socketserver.UsersDataMessage{
		Uids: Uids,
		Type: "OUT_INVITATION_DELETE",
		Data: msg,
	}

	return nil
}

func invitationToRoomResponse(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	var data socketmodels.RoomInvitationResponse
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	invitationId, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return err
	}

	messagingData := &models.UserMessagingData{}
	if err := colls.UserMessagingDataCollection.FindOne(context.Background(), bson.M{"_id": uid}).Decode(&messagingData); err != nil {
		return err
	}

	invitationIndex := -1
	for i, inv := range messagingData.Invitations {
		if inv.ID == invitationId {
			invitationIndex = i
			break
		}
	}
	if invitationIndex == -1 {
		return fmt.Errorf("Invitation not found")
	}

	authorMessagingData := &models.UserMessagingData{}
	if err := colls.UserMessagingDataCollection.FindOne(context.Background(), bson.M{"_id": messagingData.Invitations[invitationIndex].Author}).Decode(&authorMessagingData); err != nil {
		return err
	}

	roomExternalData := &models.RoomExternalData{}
	if err := colls.RoomExternalDataCollection.FindOne(context.Background(), bson.M{"_id": messagingData.Invitations[invitationIndex].RoomID}).Decode(&roomExternalData); err != nil {
		return err
	}

	var deleteIfErr error = nil
	for _, oi := range roomExternalData.Banned {
		if oi == uid {
			deleteIfErr = fmt.Errorf("You can no longer accept this invitation, you have been banned from the room")
			break
		}
	}
	for _, oi := range roomExternalData.Members {
		if oi == uid {
			deleteIfErr = fmt.Errorf("You are already a member of this room")
			break
		}
	}
	for _, oi := range authorMessagingData.Blocked {
		if oi == uid {
			deleteIfErr = fmt.Errorf("The sender of this invitation has blocked your account, invitation is no longer valid")
			break
		}
	}
	if deleteIfErr != nil {
		if err := colls.UserMessagingDataCollection.FindOneAndUpdate(context.Background(), bson.M{"_id": uid}, bson.M{
			"$pull": bson.M{
				"invitations": bson.M{
					"_id":    invitationId,
					"author": messagingData.Invitations[invitationIndex].Author,
				},
			},
		}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&messagingData); err != nil {
			return err
		}
		if !checkAnythingReceivedFrom(*messagingData, messagingData.Invitations[invitationIndex].Author) {
			if _, err := colls.UserMessagingDataCollection.UpdateByID(context.Background(), uid, bson.M{
				"$pull": bson.M{
					"messages_received_from": messagingData.Invitations[invitationIndex].Author,
				},
			}); err != nil {
				return err
			}
		}
		msg := &socketmodels.OutRoomInvitationDelete{
			ID:     invitationId.Hex(),
			Author: uid.Hex(),
		}
		Uids := make(map[primitive.ObjectID]struct{})
		Uids[uid] = struct{}{}
		Uids[messagingData.Invitations[invitationIndex].Author] = struct{}{}
		ss.SendDataToUsers <- socketserver.UsersDataMessage{
			Uids: Uids,
			Type: "OUT_INVITATION_DELETE",
			Data: msg,
		}
		return deleteIfErr
	}

	if _, err := colls.RoomExternalDataCollection.UpdateOne(context.Background(), bson.M{
		"_id":             messagingData.Invitations[invitationIndex].ID,
		"invitations._id": invitationId,
	}, bson.M{
		"$set": bson.M{
			"invitations.$.accepted": data.Accept,
			"invitations.$.declined": !data.Accept,
		},
	}); err != nil {
		return err
	}

	inv := &socketmodels.OutRoomInvitationResponse{
		ID:        invitationId.Hex(),
		Author:    uid.Hex(),
		Recipient: messagingData.Invitations[invitationIndex].Author.Hex(),
	}
	Uids := make(map[primitive.ObjectID]struct{})
	Uids[uid] = struct{}{}
	Uids[messagingData.Invitations[invitationIndex].Author] = struct{}{}
	ss.SendDataToUsers <- socketserver.UsersDataMessage{
		Uids: Uids,
		Type: "OUT_INVITATION_RESPONSE",
		Data: inv,
	}

	return nil
}

func friendRequest(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	var data socketmodels.FriendRequest
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	recipientId, err := primitive.ObjectIDFromHex(data.Recipient)
	if err != nil {
		return err
	}

	messagingData := &models.UserMessagingData{}
	if err := colls.UserMessagingDataCollection.FindOne(context.Background(), bson.M{"_id": recipientId}).Decode(&messagingData); err != nil {
		return err
	}
	for _, oi := range messagingData.Blocked {
		if oi == recipientId {
			return fmt.Errorf("You have blocked this users account. You must unblock the user before sending them a friend request")
		}
	}

	recipientMessagingData := &models.UserMessagingData{}
	if err := colls.UserMessagingDataCollection.FindOne(context.Background(), bson.M{"_id": recipientId}).Decode(&recipientMessagingData); err != nil {
		return err
	} else {
		for _, oi := range recipientMessagingData.Blocked {
			if oi == uid {
				return fmt.Errorf("This user has blocked your account, you cannot send them friend requests")
			}
		}
	}

	friendRequestId := primitive.NewObjectID()

	if _, err := colls.UserMessagingDataCollection.UpdateByID(context.Background(), recipientId, bson.M{
		"$push": bson.M{
			"friend_requests": models.FriendRequest{
				ID:        friendRequestId,
				CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
				Author:    uid,
				Accepted:  false,
				Declined:  false,
			},
		},
		"$addToSet": bson.M{
			// invitations count as messages, even though they are stored in a seperate array
			"messages_received_from": uid,
		},
	}); err != nil {
		return err
	}

	if _, err := colls.UserMessagingDataCollection.UpdateByID(context.Background(), uid, bson.M{
		"$addToSet": bson.M{
			// invitations count as messages, even though they are stored in a seperate array
			"messages_sent_to": recipientId,
		},
	}); err != nil {
		return err
	}

	msg := &socketmodels.OutFriendRequest{
		ID:        friendRequestId.Hex(),
		Author:    uid.Hex(),
		Recipient: recipientId.Hex(),
	}
	Uids := make(map[primitive.ObjectID]struct{})
	Uids[uid] = struct{}{}
	Uids[recipientId] = struct{}{}
	ss.SendDataToUsers <- socketserver.UsersDataMessage{
		Uids: Uids,
		Type: "OUT_FRIEND_REQUEST",
		Data: msg,
	}

	return nil
}

func deleteFriendRequest(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	var data socketmodels.OutFriendRequestDelete
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	friendRequestId, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return err
	}

	messagingData := &models.UserMessagingData{}
	if err := colls.UserMessagingDataCollection.FindOne(context.Background(), bson.M{"_id": uid}).Decode(&messagingData); err != nil {
		return err
	}

	friendRequestIndex := -1
	for i, fr := range messagingData.FriendRequests {
		if fr.ID == friendRequestId {
			friendRequestIndex = i
			break
		}
	}
	if friendRequestIndex == -1 {
		return fmt.Errorf("Friend request not found")
	}

	if messagingData.FriendRequests[friendRequestIndex].Author != uid {
		return fmt.Errorf("Unauthorized")
	}

	if err := colls.UserMessagingDataCollection.FindOneAndUpdate(context.Background(), bson.M{"_id": uid}, bson.M{
		"$pull": bson.M{
			"friend_requests": bson.M{
				"_id":    friendRequestIndex,
				"author": messagingData.FriendRequests[friendRequestIndex].Author,
			},
		},
	}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&messagingData); err != nil {
		return err
	} else {
		if !checkAnythingReceivedFrom(*messagingData, messagingData.FriendRequests[friendRequestIndex].Author) {
			if _, err := colls.UserMessagingDataCollection.UpdateByID(context.Background(), uid, bson.M{
				"$pull": bson.M{
					"messages_received_from": messagingData.FriendRequests[friendRequestIndex].Author,
				},
			}); err != nil {
				return err
			}
		}
	}

	msg := &socketmodels.OutRoomInvitationDelete{
		ID:     friendRequestId.Hex(),
		Author: uid.Hex(),
	}
	Uids := make(map[primitive.ObjectID]struct{})
	Uids[uid] = struct{}{}
	Uids[messagingData.FriendRequests[friendRequestIndex].Author] = struct{}{}
	ss.SendDataToUsers <- socketserver.UsersDataMessage{
		Uids: Uids,
		Type: "OUT_FRIEND_REQUEST_DELETE",
		Data: msg,
	}

	return nil
}

func friendRequestResponse(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	var data socketmodels.OutFriendRequestResponse
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	friendRequestId, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return err
	}

	messagingData := &models.UserMessagingData{}
	if err := colls.UserMessagingDataCollection.FindOne(context.Background(), bson.M{"_id": uid}).Decode(&messagingData); err != nil {
		return err
	}

	friendRequestIndex := -1
	for i, inv := range messagingData.FriendRequests {
		if inv.ID == friendRequestId {
			friendRequestIndex = i
			break
		}
	}
	if friendRequestIndex == -1 {
		return fmt.Errorf("Friend request not found")
	}

	authorMessagingData := &models.UserMessagingData{}
	if err := colls.UserMessagingDataCollection.FindOne(context.Background(), bson.M{"_id": messagingData.FriendRequests[friendRequestIndex].Author}).Decode(&authorMessagingData); err != nil {
		return err
	}

	var deleteIfErr error = nil
	for _, oi := range authorMessagingData.Blocked {
		if oi == uid {
			deleteIfErr = fmt.Errorf("The sender of this invitation has blocked your account, friend request is no longer valid")
			break
		}
	}
	if deleteIfErr != nil {
		if err := colls.UserMessagingDataCollection.FindOneAndUpdate(context.Background(), bson.M{"_id": uid}, bson.M{
			"$pull": bson.M{
				"friend_requests": bson.M{
					"_id":    friendRequestId,
					"author": messagingData.Invitations[friendRequestIndex].Author,
				},
			},
		}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&messagingData); err != nil {
			return err
		}
		if !checkAnythingReceivedFrom(*messagingData, messagingData.FriendRequests[friendRequestIndex].Author) {
			if _, err := colls.UserMessagingDataCollection.UpdateByID(context.Background(), uid, bson.M{
				"$pull": bson.M{
					"messages_received_from": messagingData.FriendRequests[friendRequestIndex].Author,
				},
			}); err != nil {
				return err
			}
		}
		msg := &socketmodels.OutFriendRequestResponse{
			ID:     friendRequestId.Hex(),
			Author: uid.Hex(),
		}
		Uids := make(map[primitive.ObjectID]struct{})
		Uids[uid] = struct{}{}
		Uids[messagingData.FriendRequests[friendRequestIndex].Author] = struct{}{}
		ss.SendDataToUsers <- socketserver.UsersDataMessage{
			Uids: Uids,
			Type: "OUT_FRIEND_REQUEST_DELETE",
			Data: msg,
		}
		return deleteIfErr
	}

	if _, err := colls.RoomExternalDataCollection.UpdateOne(context.Background(), bson.M{
		"_id":                 messagingData.FriendRequests[friendRequestIndex].ID,
		"friend_requests._id": friendRequestId,
	}, bson.M{
		"$set": bson.M{
			"friend_requests.$.accepted": data.Accept,
			"friend_requests.$.declined": !data.Accept,
		},
	}); err != nil {
		return err
	}

	inv := &socketmodels.OutFriendRequestResponse{
		ID:        friendRequestId.Hex(),
		Author:    uid.Hex(),
		Recipient: messagingData.FriendRequests[friendRequestIndex].Author.Hex(),
	}
	Uids := make(map[primitive.ObjectID]struct{})
	Uids[uid] = struct{}{}
	Uids[messagingData.FriendRequests[friendRequestIndex].Author] = struct{}{}
	ss.SendDataToUsers <- socketserver.UsersDataMessage{
		Uids: Uids,
		Type: "OUT_FRIEND_REQUEST_RESPONSE",
		Data: inv,
	}

	return nil
}

func blockUser(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	var data socketmodels.BlockUnblockUser
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	if _, err := colls.UserMessagingDataCollection.UpdateByID(context.Background(), uid, bson.M{
		"$addToSet": bson.M{
			"blocked": data.UID,
		},
		"$pull": bson.M{
			"friends": data.UID,
			"messages": bson.M{
				"author": data.UID,
			},
			"messages_sent_to":       data.UID,
			"messages_received_from": data.UID,
		},
	}); err != nil {
		return err
	}

	if _, err := colls.UserMessagingDataCollection.UpdateByID(context.Background(), data.UID, bson.M{
		"$pull": bson.M{
			"friends": uid,
			"messages": bson.M{
				"author": uid,
			},
			"messages_sent_to":       uid,
			"messages_received_from": uid,
		},
	}); err != nil {
		return nil
	}

	return nil
}

func unblockUser(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	var data socketmodels.BlockUnblockUser
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	if _, err := colls.UserMessagingDataCollection.UpdateByID(context.Background(), uid, bson.M{
		"$pull": bson.M{
			"blocked": data.UID,
		},
	}); err != nil {
		return err
	}

	return nil
}

// helper function - used to check if messages_sent_to/messages_received_from should have a uid pulled
func checkAnythingReceivedFrom(messagingData models.UserMessagingData, sender primitive.ObjectID) bool {
	for _, inv := range messagingData.Invitations {
		if inv.Author == sender {
			return true
		}
	}
	for _, msg := range messagingData.Messages {
		if msg.Author == sender {
			return true
		}
	}
	for _, fr := range messagingData.FriendRequests {
		if fr.Author == sender {
			return true
		}
	}
	return false
}
