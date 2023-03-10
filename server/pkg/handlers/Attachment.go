package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/web-stuff-98/electron-social-chat/pkg/attachmentserver"
	"github.com/web-stuff-98/electron-social-chat/pkg/db/models"
	"github.com/web-stuff-98/electron-social-chat/pkg/helpers"
	"github.com/web-stuff-98/electron-social-chat/pkg/socketserver"
	"github.com/web-stuff-98/electron-social-chat/pkg/validation"
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
				responseMessage(w, http.StatusNotFound, "Channel not found")
			} else {
				responseMessage(w, http.StatusInternalServerError, "Internal error")
			}
			return
		}
		channelMessages := &models.RoomChannelMessages{}
		if err := h.Collections.RoomChannelMessagesCollection.FindOne(r.Context(), bson.M{"_id": recipient}).Decode(&channelMessages); err != nil {
			if err == mongo.ErrNoDocuments {
				responseMessage(w, http.StatusNotFound, "Channel not found")
			} else {
				responseMessage(w, http.StatusInternalServerError, "Internal error")
			}
			return
		}
		roomExternalData := &models.RoomExternalData{}
		if err := h.Collections.RoomExternalDataCollection.FindOne(r.Context(), bson.M{"_id": channel.RoomID}).Decode(&roomExternalData); err != nil {
			if err == mongo.ErrNoDocuments {
				responseMessage(w, http.StatusNotFound, "Room not found")
			} else {
				responseMessage(w, http.StatusInternalServerError, "Internal error")
			}
			return
		}
		found := false
		for _, rcm := range channelMessages.Messages {
			if rcm.ID == msgId {
				found = true
				if rcm.Author != user.ID {
					responseMessage(w, http.StatusUnauthorized, "Unauthorized")
					return
				}
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
				if dm.Author != user.ID {
					responseMessage(w, http.StatusUnauthorized, "Unauthorized")
					return
				}
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
	if err != nil {
		responseMessage(w, http.StatusBadRequest, "Error reading body")
		return
	}

	if len(body) > 4*1024*1024 {
		responseMessage(w, http.StatusRequestEntityTooLarge, "Chunk too large")
		return
	}

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
	log.Println("CHUNK BEING CREATED")
	recvChan := make(chan bool)
	h.AttachmentServer.ChunkChan <- attachmentserver.InChunk{
		Uid:           user.ID,
		MsgId:         msgId,
		Data:          body,
		SendUpdatesTo: sendUpdatesTo,
		RecvChan:      recvChan,
	}
	success := <-recvChan
	if success {
		log.Println("CHUNK CREATED")
		responseMessage(w, http.StatusOK, "Chunk created")
	} else {
		log.Println("CHUNK FAILED")
		responseMessage(w, http.StatusInternalServerError, "Chunk failed")
	}
}

func (h handler) CreateAttachmentMetadata(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetUserFromRequest(r, r.Context(), *h.Collections, h.RedisClient)
	if err != nil {
		responseMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseMessage(w, http.StatusBadRequest, "Bad request")
		return
	}
	var dataInput validation.AttachmentMetadata
	if err := json.Unmarshal(body, &dataInput); err != nil {
		responseMessage(w, http.StatusBadRequest, "Bad request")
		return
	}
	validate := validator.New()
	if err := validate.Struct(dataInput); err != nil {
		responseMessage(w, http.StatusBadRequest, "Bad request")
		return
	}

	msgId, err := primitive.ObjectIDFromHex(dataInput.ID)
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
				responseMessage(w, http.StatusNotFound, "Channel not found")
			} else {
				responseMessage(w, http.StatusInternalServerError, "Internal error")
			}
			return
		}
		channelMessages := &models.RoomChannelMessages{}
		if err := h.Collections.RoomChannelMessagesCollection.FindOne(r.Context(), bson.M{"_id": recipient}).Decode(&channelMessages); err != nil {
			if err == mongo.ErrNoDocuments {
				responseMessage(w, http.StatusNotFound, "Channel not found")
			} else {
				responseMessage(w, http.StatusInternalServerError, "Internal error")
			}
			return
		}
		roomExternalData := &models.RoomExternalData{}
		if err := h.Collections.RoomExternalDataCollection.FindOne(r.Context(), bson.M{"_id": channel.RoomID}).Decode(&roomExternalData); err != nil {
			if err == mongo.ErrNoDocuments {
				responseMessage(w, http.StatusNotFound, "Room not found")
			} else {
				responseMessage(w, http.StatusInternalServerError, "Internal error")
			}
			return
		}
		found := false
		for _, rcm := range channelMessages.Messages {
			if rcm.ID == msgId {
				found = true
				if rcm.Author != user.ID {
					responseMessage(w, http.StatusUnauthorized, "Unauthorized")
					return
				}
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
				if dm.Author != user.ID {
					responseMessage(w, http.StatusUnauthorized, "Unauthorized")
					return
				}
				break
			}
		}
		if !found {
			responseMessage(w, http.StatusNotFound, "Message not found")
			return
		}
	}

	if dataInput.Size > 20*1024*1024 {
		if _, err := h.Collections.AttachmentMetadataCollection.InsertOne(r.Context(), models.AttachmentData{
			ID:     msgId,
			Meta:   dataInput.MimeType,
			Name:   dataInput.Name,
			Size:   dataInput.Size,
			Ratio:  0,
			Failed: true,
		}); err != nil {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
			return
		}
		responseMessage(w, http.StatusRequestEntityTooLarge, "File too large")
		return
	} else {
		if _, err := h.Collections.AttachmentMetadataCollection.InsertOne(r.Context(), models.AttachmentData{
			ID:     msgId,
			Meta:   dataInput.MimeType,
			Name:   dataInput.Name,
			Size:   dataInput.Size,
			Ratio:  0,
			Failed: false,
		}); err != nil {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
			return
		}
	}

	responseMessage(w, http.StatusCreated, "Metadata created")
}

// Download attachment as a file using octet stream
func (h handler) DownloadAttachment(w http.ResponseWriter, r *http.Request) {
	rawMsgId := mux.Vars(r)["id"]
	msgId, err := primitive.ObjectIDFromHex(rawMsgId)
	if err != nil {
		responseMessage(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var metaData models.AttachmentData
	if h.Collections.AttachmentMetadataCollection.FindOne(r.Context(), bson.M{"_id": msgId}).Decode(&metaData); err != nil {
		if err == mongo.ErrNoDocuments {
			responseMessage(w, http.StatusNotFound, "Not found")
		} else {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
		}
		return
	}

	if metaData.Failed {
		responseMessage(w, http.StatusBadRequest, "This attachment failed to upload correctly")
		return
	}

	if metaData.Ratio != float32(1) {
		responseMessage(w, http.StatusBadRequest, "This attachment is not yet complete")
		return
	}

	var firstChunk models.AttachmentChunk
	if err := h.Collections.AttachmentChunkCollection.FindOne(r.Context(), bson.M{"_id": msgId}).Decode(&firstChunk); err != nil {
		if err == mongo.ErrNoDocuments {
			responseMessage(w, http.StatusNotFound, "Not found")
		} else {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
		}
		return
	}

	w.Header().Add("Content-Type", "application/octet-stream")
	w.Header().Add("Content-Length", strconv.Itoa(metaData.Size))
	w.Header().Add("Content-Disposition", `attachment; filename="`+metaData.Name+`"`)

	w.Write(firstChunk.Data.Data)

	if firstChunk.NextChunkID != primitive.NilObjectID {
		recursivelyWriteAttachmentChunksToResponse(w, firstChunk.NextChunkID, h.Collections.AttachmentChunkCollection, r.Context())
	}
}

func recursivelyWriteAttachmentChunksToResponse(w http.ResponseWriter, NextChunkID primitive.ObjectID, chunkColl *mongo.Collection, ctx context.Context) error {
	var chunk models.AttachmentChunk
	if err := chunkColl.FindOne(ctx, bson.M{"_id": NextChunkID}).Decode(&chunk); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		} else {
			return err
		}
	} else {
		w.Write(chunk.Data.Data)
		if chunk.NextChunkID != primitive.NilObjectID {
			return recursivelyWriteAttachmentChunksToResponse(w, chunk.NextChunkID, chunkColl, ctx)
		} else {
			return nil
		}
	}
}
