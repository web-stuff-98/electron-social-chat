package handlers

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/web-stuff-98/electron-social-chat/pkg/attachmentserver"
	"github.com/web-stuff-98/electron-social-chat/pkg/db/models"
	"github.com/web-stuff-98/electron-social-chat/pkg/helpers"
	"github.com/web-stuff-98/electron-social-chat/pkg/socketserver"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h handler) UploadAttachmentChunk(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetUserFromRequest(r, r.Context(), *h.Collections, h.RedisClient)
	if err != nil {
		responseMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	msgId, err := primitive.ObjectIDFromHex(mux.Vars(r)["msgId"])
	if err != nil {
		responseMessage(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var recipient primitive.ObjectID
	isRoomMsg := r.URL.Query().Has("channel_id")
	if isRoomMsg {
		if channelId, err := primitive.ObjectIDFromHex(r.URL.Query().Get("channel_id")); err != nil {
			responseMessage(w, http.StatusBadRequest, "Invalid ID")
			return
		} else {
			recipient = channelId
		}
		channel := &models.RoomChannel{}
		if err := h.Collections.RoomChannelCollection.FindOne(r.Context(), bson.M{"_id": recipient}).Decode(&channel); err != nil {
			if err == mongo.ErrNoDocuments {
				responseMessage(w, http.StatusNotFound, "User not found")
			} else {
				responseMessage(w, http.StatusInternalServerError, "Internal error")
			}
			return
		}
		roomExternalData := &models.RoomExternalData{}
		if err := h.Collections.RoomExternalDataCollection.FindOne(r.Context(), bson.M{"_id": recipient}).Decode(&roomExternalData); err != nil {
			if err == mongo.ErrNoDocuments {
				responseMessage(w, http.StatusNotFound, "Room not found")
			} else {
				responseMessage(w, http.StatusInternalServerError, "Internal error")
			}
			return
		}
		found := false
		for _, rcm := range channel.Messages {
			if rcm.ID == msgId {
				found = true
				break
			}
		}
		if !found {
			responseMessage(w, http.StatusNotFound, "Message not found")
			return
		}
	} else {
		if uid, err := primitive.ObjectIDFromHex(r.URL.Query().Get("uid")); err != nil {
			responseMessage(w, http.StatusBadRequest, "Invalid ID")
			return
		} else {
			recipient = uid
		}
		recipientMessagingData := &models.UserMessagingData{}
		if err := h.Collections.UserMessagingDataCollection.FindOne(r.Context(), bson.M{"_id": recipient}).Decode(&recipientMessagingData); err != nil {
			if err == mongo.ErrNoDocuments {
				responseMessage(w, http.StatusNotFound, "User not found")
			} else {
				responseMessage(w, http.StatusInternalServerError, "Internal error")
			}
			return
		}
		found := false
		for _, dm := range recipientMessagingData.Messages {
			if dm.ID == msgId {
				found = true
				break
			}
		}
		if !found {
			responseMessage(w, http.StatusNotFound, "Message not found")
			return
		}
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	sendUpdatesTo := make(map[primitive.ObjectID]struct{})
	if isRoomMsg {
		recvChan := make(chan map[primitive.ObjectID]struct{})
		h.SocketServer.GetSubscriptionUids <- socketserver.GetSubscriptionUids{
			Name:     "channel:" + recipient.Hex(),
			RecvChan: recvChan,
		}
		uids := <-recvChan
		sendUpdatesTo = uids
	} else {
		sendUpdatesTo[user.ID] = struct{}{}
		sendUpdatesTo[recipient] = struct{}{}
	}
	h.AttachmentServer.ChunkChan <- attachmentserver.InChunk{
		Uid:           user.ID,
		MsgId:         msgId,
		Data:          body,
		SendUpdatesTo: sendUpdatesTo,
	}

	responseMessage(w, http.StatusOK, "Chunk created")
}
