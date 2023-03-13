package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/web-stuff-98/electron-social-chat/pkg/attachmentserver"
	"github.com/web-stuff-98/electron-social-chat/pkg/db"
	"github.com/web-stuff-98/electron-social-chat/pkg/db/models"
	"github.com/web-stuff-98/electron-social-chat/pkg/socketmodels"
	"github.com/web-stuff-98/electron-social-chat/pkg/socketserver"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func HandleSocketEvent(eventType string, data []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, as *attachmentserver.AttachmentServer, colls *db.Collections) error {
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
		err := roomMessageDelete(data, conn, uid, ss, as, colls)
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
	case "FRIEND_REQUEST_DELETE":
		err := deleteFriendRequest(data, conn, uid, ss, colls)
		return err
	case "FRIEND_REQUEST_RESPONSE":
		err := friendRequestResponse(data, conn, uid, ss, colls)
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
	case "BLOCK":
		err := blockUser(data, conn, uid, ss, as, colls)
		return err
	case "UNBLOCK":
		err := unblockUser(data, conn, uid, ss, as, colls)
		return err
	case "BAN":
		err := banUser(data, conn, uid, ss, as, colls)
		return err
	case "UNBAN":
		err := unbanUser(data, conn, uid, ss, as, colls)
		return err
	case "CALL_USER":
		err := callUser(data, conn, uid, ss, colls)
		return err
	case "CALL_USER_RESPONSE":
		err := callUserResponse(data, conn, uid, ss, colls)
		return err
	case "CALL_LEAVE":
		err := callLeave(data, conn, uid, ss, colls)
		return err
	case "CALL_WEBRTC_OFFER":
		err := callWebRTCOffer(data, conn, uid, ss, colls)
		return err
	case "CALL_WEBRTC_ANSWER":
		err := callWebRTCAnswer(data, conn, uid, ss, colls)
		return err
	}
	return fmt.Errorf("Unrecognized event type:", eventType)
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
				ID:            msgId,
				Content:       data.Content,
				CreatedAt:     primitive.NewDateTimeFromTime(time.Now()),
				UpdatedAt:     primitive.NewDateTimeFromTime(time.Now()),
				Author:        uid,
				HasAttachment: data.HasAttachment,
			},
		},
	}); err != nil {
		return err
	}

	if outBytes, err := json.Marshal(socketmodels.OutRoomMessage{
		Type:          "OUT_ROOM_MESSAGE",
		Content:       data.Content,
		ID:            msgId.Hex(),
		Author:        uid.Hex(),
		HasAttachment: data.HasAttachment,
	}); err == nil {
		ss.SendDataToSubscription <- socketserver.SubscriptionDataMessage{
			Name: "channel:" + channelId.Hex(),
			Data: outBytes,
		}
	}

	if data.HasAttachment {
		ss.SendDataToUser <- socketserver.UserDataMessage{
			Type: "ATTACHMENT_REQUEST",
			Uid:  uid,
			Data: socketmodels.AttachmentRequest{
				MsgID:  msgId.Hex(),
				IsRoom: true,
			},
		}
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

func roomMessageDelete(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, as *attachmentserver.AttachmentServer, colls *db.Collections) error {
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

	as.DeleteChan <- attachmentserver.Delete{
		MsgId: msgId,
		Uid:   uid,
	}

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
				ID:            msgId,
				CreatedAt:     primitive.NewDateTimeFromTime(time.Now()),
				UpdatedAt:     primitive.NewDateTimeFromTime(time.Now()),
				Author:        uid,
				Content:       data.Content,
				HasAttachment: data.HasAttachment,
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

	Uids := make(map[primitive.ObjectID]struct{})
	Uids[uid] = struct{}{}
	Uids[recipientId] = struct{}{}
	ss.SendDataToUsers <- socketserver.UsersDataMessage{
		Uids: Uids,
		Type: "OUT_DIRECT_MESSAGE",
		Data: socketmodels.OutDirectMessage{
			ID:            msgId.Hex(),
			Content:       data.Content,
			Author:        uid.Hex(),
			Recipient:     recipientId.Hex(),
			HasAttachment: data.HasAttachment,
		},
	}

	if data.HasAttachment {
		ss.SendDataToUser <- socketserver.UserDataMessage{
			Type: "ATTACHMENT_REQUEST",
			Uid:  uid,
			Data: socketmodels.AttachmentRequest{
				MsgID:  msgId.Hex(),
				IsRoom: false,
			},
		}
	}

	return nil
}

func directMessageUpdate(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	var data socketmodels.DirectMessageUpdate
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	if strings.TrimSpace(data.Content) == "" {
		return fmt.Errorf("Cannot submit an empty message")
	}

	if len(data.Content) > 300 {
		return fmt.Errorf("Max 300 characters")
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
	Uids := make(map[primitive.ObjectID]struct{})
	Uids[uid] = struct{}{}
	Uids[recipientId] = struct{}{}
	ss.SendDataToUsers <- socketserver.UsersDataMessage{
		Uids: Uids,
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

	if recipientId == uid {
		return fmt.Errorf("You cannot send an invitation to yourself")
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
		Type: "OUT_ROOM_INVITATION",
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
		ID:        invitationId.Hex(),
		Author:    uid.Hex(),
		Recipient: data.Recipient,
	}
	Uids := make(map[primitive.ObjectID]struct{})
	Uids[uid] = struct{}{}
	Uids[messagingData.Invitations[invitationIndex].Author] = struct{}{}
	ss.SendDataToUsers <- socketserver.UsersDataMessage{
		Uids: Uids,
		Type: "OUT_ROOM_INVITATION_DELETE",
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
	roomInternalData := &models.RoomInternalData{}
	if err := colls.RoomInternalDataCollection.FindOne(context.Background(), bson.M{"_id": messagingData.Invitations[invitationIndex].RoomID}).Decode(&roomInternalData); err != nil {
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
			ID:        invitationId.Hex(),
			Author:    messagingData.Invitations[invitationIndex].Author.Hex(),
			Recipient: uid.Hex(),
		}
		Uids := make(map[primitive.ObjectID]struct{})
		Uids[uid] = struct{}{}
		Uids[messagingData.Invitations[invitationIndex].Author] = struct{}{}
		ss.SendDataToUsers <- socketserver.UsersDataMessage{
			Uids: Uids,
			Type: "OUT_ROOM_INVITATION_DELETE",
			Data: msg,
		}
		return deleteIfErr
	}

	if _, err := colls.UserMessagingDataCollection.UpdateOne(context.Background(), bson.M{
		"_id":             messagingData.ID,
		"invitations._id": invitationId,
	}, bson.M{
		"$set": bson.M{
			"invitations.$.accepted": data.Accept,
			"invitations.$.declined": !data.Accept,
		},
	}); err != nil {
		return err
	}

	if data.Accept {
		if _, err := colls.RoomExternalDataCollection.UpdateOne(context.Background(), bson.M{
			"_id": messagingData.Invitations[invitationIndex].RoomID,
		}, bson.M{
			"$addToSet": bson.M{
				"members": messagingData.ID,
			},
		}); err != nil {
			return err
		}
		for _, channelId := range roomInternalData.Channels {
			recvChan := make(chan map[primitive.ObjectID]struct{})
			ss.GetSubscriptionUids <- socketserver.GetSubscriptionUids{
				RecvChan: recvChan,
				Name:     "channel:" + channelId.Hex(),
			}
			uidsInChannel := <-recvChan
			ss.SendDataToUsers <- socketserver.UsersDataMessage{
				Uids: uidsInChannel,
				Type: "MEMBER_ADDED",
				Data: socketmodels.MemberAdded{
					Uid:    uid.Hex(),
					RoomID: messagingData.Invitations[invitationIndex].RoomID.Hex(),
				},
			}
		}
	}

	Uids := make(map[primitive.ObjectID]struct{})
	Uids[uid] = struct{}{}
	Uids[messagingData.Invitations[invitationIndex].Author] = struct{}{}
	ss.SendDataToUsers <- socketserver.UsersDataMessage{
		Uids: Uids,
		Type: "OUT_ROOM_INVITATION_RESPONSE",
		Data: socketmodels.OutRoomInvitationResponse{
			ID:        invitationId.Hex(),
			Author:    messagingData.Invitations[invitationIndex].Author.Hex(),
			Recipient: uid.Hex(),
			Accept:    data.Accept,
		},
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

	if recipientId == uid {
		return fmt.Errorf("You cannot send a friend request to yourself")
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

	for _, fr := range messagingData.FriendRequests {
		if fr.Author == recipientId {
			return fmt.Errorf("This user has already sent you a friend request")
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
		for _, fr := range recipientMessagingData.FriendRequests {
			if fr.Author == uid && !fr.Accepted && !fr.Declined {
				return fmt.Errorf("You have already sent this user a friend request")
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
			// friend requests count as messages, even though they are stored in a seperate array
			"messages_received_from": uid,
		},
	}); err != nil {
		return err
	}

	if _, err := colls.UserMessagingDataCollection.UpdateByID(context.Background(), uid, bson.M{
		"$addToSet": bson.M{
			// friend requests as messages, even though they are stored in a seperate array
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

	msg := &socketmodels.OutFriendRequestDelete{
		ID:        friendRequestId.Hex(),
		Author:    uid.Hex(),
		Recipient: data.Recipient,
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
	var data socketmodels.FriendRequestResponse
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

	frq := messagingData.FriendRequests[friendRequestIndex]

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
			ID:        friendRequestId.Hex(),
			Author:    uid.Hex(),
			Recipient: data.Recipient,
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

	if data.Accept {
		if _, err := colls.UserMessagingDataCollection.UpdateOne(context.Background(), bson.M{
			"_id":                 uid,
			"friend_requests._id": friendRequestId,
		}, bson.M{
			"$set": bson.M{
				"friend_requests.$.accepted": data.Accept,
				"friend_requests.$.declined": !data.Accept,
			},
			"$addToSet": bson.M{
				"friends": frq.Author,
			},
		}); err != nil {
			return err
		}

		if _, err := colls.UserMessagingDataCollection.UpdateOne(context.Background(), bson.M{
			"_id": frq.Author,
		}, bson.M{
			"$addToSet": bson.M{
				"friends": uid,
			},
		}); err != nil {
			return err
		}
	}

	inv := &socketmodels.OutFriendRequestResponse{
		ID:        friendRequestId.Hex(),
		Author:    uid.Hex(),
		Recipient: frq.Author.Hex(),
		Accept:    data.Accept,
	}
	Uids := make(map[primitive.ObjectID]struct{})
	Uids[uid] = struct{}{}
	Uids[frq.Author] = struct{}{}
	ss.SendDataToUsers <- socketserver.UsersDataMessage{
		Uids: Uids,
		Type: "OUT_FRIEND_REQUEST_RESPONSE",
		Data: inv,
	}

	return nil
}

func blockUser(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, as *attachmentserver.AttachmentServer, colls *db.Collections) error {
	var data socketmodels.Block
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	blockedUid, err := primitive.ObjectIDFromHex(data.Uid)
	if err != nil {
		return err
	}

	userMessagingData := &models.UserMessagingData{}
	if err := colls.UserMessagingDataCollection.FindOneAndUpdate(context.Background(), bson.M{"_id": uid}, bson.M{
		"$addToSet": bson.M{
			"blocked": blockedUid,
		},
		"$pull": bson.M{
			"friends": blockedUid,
			"messages": bson.M{
				"author": blockedUid,
			},
			"messages_sent_to":       blockedUid,
			"messages_received_from": blockedUid,
		},
	}, options.FindOneAndUpdate().SetReturnDocument(options.Before)).Decode(&userMessagingData); err != nil {
		return err
	}
	for _, dm := range userMessagingData.Messages {
		if dm.Author == blockedUid && dm.HasAttachment {
			as.DeleteChan <- attachmentserver.Delete{
				MsgId: dm.ID,
				Uid:   blockedUid,
			}
		}
	}

	blockedUserMessagingData := &models.UserMessagingData{}
	if err := colls.UserMessagingDataCollection.FindOneAndUpdate(context.Background(), bson.M{"_id": blockedUid}, bson.M{
		"$pull": bson.M{
			"friends": uid,
			"messages": bson.M{
				"author": uid,
			},
			"messages_sent_to":       uid,
			"messages_received_from": uid,
		},
	}, options.FindOneAndUpdate().SetReturnDocument(options.Before)).Decode(&blockedUserMessagingData); err != nil {
		return err
	}
	for _, dm := range blockedUserMessagingData.Messages {
		if dm.Author == uid && dm.HasAttachment {
			as.DeleteChan <- attachmentserver.Delete{
				MsgId: dm.ID,
				Uid:   uid,
			}
		}
	}

	rooms := []models.Room{}
	if cursor, err := colls.RoomCollection.Find(context.Background(), bson.M{"author": uid}); err != nil {
		cursor.Close(context.Background())
		return err
	} else {
		cursor.All(context.Background(), rooms)
		defer cursor.Close(context.Background())
		roomIds := []primitive.ObjectID{}
		for _, r := range rooms {
			roomIds = append(roomIds, r.ID)
			internalData := &models.RoomInternalData{}
			if err := colls.RoomInternalDataCollection.FindOne(context.Background(), bson.M{"_id": r.ID}).Decode(&internalData); err != nil {
				return err
			}
			for _, oi := range internalData.Channels {
				recvChan := make(chan map[primitive.ObjectID]struct{})
				ss.GetSubscriptionUids <- socketserver.GetSubscriptionUids{
					RecvChan: recvChan,
					Name:     "channel:" + oi.Hex(),
				}
				uidsInChannel := <-recvChan
				// Blocking a user also bans them from all the blockers rooms
				ss.SendDataToUsers <- socketserver.UsersDataMessage{
					Uids: uidsInChannel,
					Data: socketmodels.Banned{
						Banner: uid.Hex(),
						Banned: data.Uid,
						RoomID: r.ID.Hex(),
					},
					Type: "BANNED",
				}
				ss.RemoveUserFromSubscription <- socketserver.RemoveUserFromSubscription{
					Name: "channel:" + oi.Hex(),
					Uid:  blockedUid,
				}
				channelMessages := &models.RoomChannelMessages{}
				if err := colls.RoomChannelMessagesCollection.FindOneAndUpdate(context.Background(), bson.M{"_id": oi}, bson.M{
					"messages": bson.M{
						"author": blockedUid,
					},
				}, options.FindOneAndUpdate().SetReturnDocument(options.Before)).Decode(&channelMessages); err != nil {
					return err
				}
				for _, rcm := range channelMessages.Messages {
					if rcm.Author == blockedUid && rcm.HasAttachment {
						as.DeleteChan <- attachmentserver.Delete{
							MsgId: rcm.ID,
							Uid:   blockedUid,
						}
					}
				}
			}
		}
		if _, err := colls.RoomExternalDataCollection.UpdateMany(context.Background(), bson.M{"_id": bson.M{"$in": roomIds}}, bson.M{
			"$pull": bson.M{
				"members": blockedUid,
			},
			"$addToSet": bson.M{
				"banned": blockedUid,
			},
		}); err != nil {
			return err
		}
	}

	uids := make(map[primitive.ObjectID]struct{})
	uids[uid] = struct{}{}
	uids[blockedUid] = struct{}{}

	ss.SendDataToUsers <- socketserver.UsersDataMessage{
		Uids: uids,
		Data: socketmodels.Blocked{
			Blocker: uid.Hex(),
		},
		Type: "BLOCKED",
	}

	return nil
}

func unblockUser(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, as *attachmentserver.AttachmentServer, colls *db.Collections) error {
	var data socketmodels.Block
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	blockedUid, err := primitive.ObjectIDFromHex(data.Uid)
	if err != nil {
		return err
	}

	if _, err := colls.UserMessagingDataCollection.UpdateByID(context.Background(), uid, bson.M{
		"$pull": bson.M{
			"blocked": blockedUid,
		},
	}); err != nil {
		return err
	}

	uids := make(map[primitive.ObjectID]struct{})
	uids[uid] = struct{}{}
	uids[blockedUid] = struct{}{}

	ss.SendDataToUsers <- socketserver.UsersDataMessage{
		Uids: uids,
		Data: socketmodels.Blocked{
			Blocker: uid.Hex(),
		},
		Type: "UNBLOCKED",
	}

	return nil
}

func banUser(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, as *attachmentserver.AttachmentServer, colls *db.Collections) error {
	var data socketmodels.Ban
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	bannedUid, err := primitive.ObjectIDFromHex(data.Uid)
	if err != nil {
		return err
	}
	roomId, err := primitive.ObjectIDFromHex(data.RoomID)
	if err != nil {
		return err
	}

	if _, err := colls.RoomExternalDataCollection.UpdateByID(context.Background(), roomId, bson.M{
		"$addToSet": bson.M{
			"banned": bannedUid,
		},
		"$pull": bson.M{
			"members": bannedUid,
		},
	}); err != nil {
		return err
	}

	internalData := &models.RoomInternalData{}
	if err := colls.RoomInternalDataCollection.FindOne(context.Background(), bson.M{"_id": roomId}).Decode(&internalData); err != nil {
		return err
	}

	uids := make(map[primitive.ObjectID]struct{})
	uids[bannedUid] = struct{}{}
	uids[uid] = struct{}{}

	for _, oi := range internalData.Channels {
		channelMessages := &models.RoomChannelMessages{}
		if err := colls.RoomChannelMessagesCollection.FindOneAndUpdate(context.Background(), bson.M{"_id": oi}, bson.M{
			"$pull": bson.M{
				"messages": bson.M{
					"author": bannedUid,
				},
			},
		}, options.FindOneAndUpdate().SetReturnDocument(options.Before)).Decode(&channelMessages); err != nil {
			return err
		}
		for _, rcm := range channelMessages.Messages {
			if rcm.Author == bannedUid && rcm.HasAttachment {
				as.DeleteChan <- attachmentserver.Delete{
					MsgId: rcm.ID,
					Uid:   bannedUid,
				}
			}
		}
		recvChan := make(chan map[primitive.ObjectID]struct{})
		ss.GetSubscriptionUids <- socketserver.GetSubscriptionUids{
			RecvChan: recvChan,
			Name:     "channel:" + oi.Hex(),
		}
		uidsInChannel := <-recvChan
		for oi2 := range uidsInChannel {
			uids[oi2] = struct{}{}
		}
		ss.RemoveUserFromSubscription <- socketserver.RemoveUserFromSubscription{
			Name: "channel:" + oi.Hex(),
			Uid:  bannedUid,
		}
	}

	ss.SendDataToUsers <- socketserver.UsersDataMessage{
		Uids: uids,
		Data: socketmodels.Banned{
			Banner: uid.Hex(),
			Banned: data.Uid,
			RoomID: data.RoomID,
		},
		Type: "BANNED",
	}

	return nil
}

func unbanUser(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, as *attachmentserver.AttachmentServer, colls *db.Collections) error {
	var data socketmodels.Ban
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	bannedUid, err := primitive.ObjectIDFromHex(data.Uid)
	if err != nil {
		return err
	}
	roomId, err := primitive.ObjectIDFromHex(data.RoomID)
	if err != nil {
		return err
	}

	if _, err := colls.RoomExternalDataCollection.UpdateByID(context.Background(), roomId, bson.M{
		"$pull": bson.M{
			"banned": bannedUid,
		},
	}); err != nil {
		return err
	}

	internalData := &models.RoomInternalData{}
	if err := colls.RoomInternalDataCollection.FindOne(context.Background(), bson.M{"_id": roomId}).Decode(&internalData); err != nil {
		return err
	}

	uids := make(map[primitive.ObjectID]struct{})
	uids[bannedUid] = struct{}{}
	uids[uid] = struct{}{}

	for _, oi := range internalData.Channels {
		channelMessages := &models.RoomChannelMessages{}
		if err := colls.RoomChannelMessagesCollection.FindOne(context.Background(), bson.M{"_id": oi}).Decode(&channelMessages); err != nil {
			return err
		}
		recvChan := make(chan map[primitive.ObjectID]struct{})
		ss.GetSubscriptionUids <- socketserver.GetSubscriptionUids{
			RecvChan: recvChan,
			Name:     "channel:" + oi.Hex(),
		}
		uidsInChannel := <-recvChan
		for oi2 := range uidsInChannel {
			uids[oi2] = struct{}{}
		}
	}

	ss.SendDataToUsers <- socketserver.UsersDataMessage{
		Uids: uids,
		Data: socketmodels.Banned{
			Banner: uid.Hex(),
			Banned: data.Uid,
		},
		Type: "UNBANNED",
	}

	return nil
}

func callUser(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	var data socketmodels.CallUser
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	callUid, err := primitive.ObjectIDFromHex(data.Uid)
	if err != nil {
		return err
	}

	calledMessagingData := &models.UserMessagingData{}
	if err := colls.UserMessagingDataCollection.FindOne(context.Background(), bson.M{"_id": callUid}).Decode(&calledMessagingData); err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("User not found")
		}
		return fmt.Errorf("Internal server error")
	}

	for _, oi := range calledMessagingData.Blocked {
		if oi == uid {
			return fmt.Errorf("This user blocked your account")
		}
	}
	foundFriend := false
	for _, oi := range calledMessagingData.Friends {
		if oi == uid {
			foundFriend = true
			break
		}
	}
	if !foundFriend {
		return fmt.Errorf("You can only call users you are friends with")
	}

	ss.CallsPendingChan <- socketserver.InCall{
		Caller: uid,
		Called: callUid,
	}

	return nil
}

func callUserResponse(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	var data socketmodels.CallResponse
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	callerUid, err := primitive.ObjectIDFromHex(data.Caller)
	if err != nil {
		return err
	}
	calledUid, err := primitive.ObjectIDFromHex(data.Called)
	if err != nil {
		return err
	}
	if callerUid == uid && data.Accept {
		return fmt.Errorf("You cannot accept a call to a another user on your own behalf")
	}

	ss.ResponseToCallChan <- socketserver.InCallResponse{
		Caller: callerUid,
		Called: calledUid,
		Accept: data.Accept,
	}

	return nil
}

func callLeave(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	var data socketmodels.CallLeave
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	ss.LeaveCallChan <- uid

	return nil
}

func callWebRTCOffer(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	var data socketmodels.CallWebRTCOfferAnswer
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	ss.SendCallRecipientOffer <- socketserver.CallerSignal{
		Signal: data.Signal,
		Caller: uid,
	}

	return nil
}

func callWebRTCAnswer(b []byte, conn *websocket.Conn, uid primitive.ObjectID, ss *socketserver.SocketServer, colls *db.Collections) error {
	var data socketmodels.CallWebRTCOfferAnswer
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	ss.SendCalledAnswer <- socketserver.CalledSignal{
		Signal: data.Signal,
		Called: uid,
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
