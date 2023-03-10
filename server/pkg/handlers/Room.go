package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
	"github.com/web-stuff-98/electron-social-chat/pkg/db/models"
	"github.com/web-stuff-98/electron-social-chat/pkg/helpers"
	"github.com/web-stuff-98/electron-social-chat/pkg/socketmodels"
	"github.com/web-stuff-98/electron-social-chat/pkg/socketserver"
	"github.com/web-stuff-98/electron-social-chat/pkg/validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetUserFromRequest(r, r.Context(), *h.Collections, h.RedisClient)
	if err != nil {
		responseMessage(w, http.StatusUnauthorized, "StatusUnauthorized")
		return
	}

	if cursor, err := h.Collections.RoomCollection.Find(r.Context(), bson.M{"author": user.ID}); err != nil {
		cursor.Close(r.Context())
		if err != mongo.ErrNoDocuments {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
			return
		}
	} else {
		var count int8
		for cursor.Next(r.Context()) {
			count++
		}
		cursor.Close(r.Context())
		if count >= 4 {
			responseMessage(w, http.StatusBadRequest, "Too many rooms. Max 4.")
			return
		}
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	var roomInput validation.Room
	if json.Unmarshal(body, &roomInput); err != nil {
		responseMessage(w, http.StatusBadRequest, "Bad request")
		return
	}
	validate := validator.New()
	if err := validate.Struct(roomInput); err != nil {
		responseMessage(w, http.StatusBadRequest, "Bad request")
		return
	}

	res := h.Collections.RoomCollection.FindOne(r.Context(), bson.M{
		"name": bson.M{
			"$regex":   roomInput.Name,
			"$options": "i",
		},
		"author": user.ID,
	})
	if res.Err() != nil {
		if res.Err() != mongo.ErrNoDocuments {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
			return
		}
	} else {
		responseMessage(w, http.StatusBadRequest, "You already have a room by that name")
		return
	}

	inserted, err := h.Collections.RoomCollection.InsertOne(r.Context(), models.Room{
		ID:     primitive.NewObjectID(),
		Name:   roomInput.Name,
		Author: user.ID,
	})
	if err != nil {
		responseMessage(w, http.StatusInternalServerError, "Internal error")
		return
	}
	mainChannelId := primitive.NewObjectID()
	mainChannel := models.RoomChannel{
		ID:          mainChannelId,
		RoomID:      inserted.InsertedID.(primitive.ObjectID),
		Name:        "Main channel",
		ToBeDeleted: false,
	}
	if _, err := h.Collections.RoomInternalDataCollection.InsertOne(r.Context(), models.RoomInternalData{
		ID:          inserted.InsertedID.(primitive.ObjectID),
		Channels:    []primitive.ObjectID{mainChannel.ID},
		MainChannel: mainChannelId,
	}); err != nil {
		h.Collections.RoomCollection.DeleteOne(r.Context(), bson.M{"_id": inserted.InsertedID.(primitive.ObjectID)})
		responseMessage(w, http.StatusInternalServerError, "Internal error")
		return
	}
	if _, err := h.Collections.RoomExternalDataCollection.InsertOne(r.Context(), models.RoomExternalData{
		ID:      inserted.InsertedID.(primitive.ObjectID),
		Private: roomInput.Private,
		Members: []primitive.ObjectID{},
		Banned:  []primitive.ObjectID{},
	}); err != nil {
		h.Collections.RoomCollection.DeleteOne(r.Context(), bson.M{"_id": inserted.InsertedID.(primitive.ObjectID)})
		responseMessage(w, http.StatusInternalServerError, "Internal error")
		return
	}
	if _, err := h.Collections.RoomChannelCollection.InsertOne(r.Context(), mainChannel); err != nil {
		h.Collections.RoomCollection.DeleteOne(r.Context(), bson.M{"_id": inserted.InsertedID.(primitive.ObjectID)})
		responseMessage(w, http.StatusInternalServerError, "Internal error")
		return
	}
	if _, err := h.Collections.RoomChannelMessagesCollection.InsertOne(r.Context(), models.RoomChannelMessages{
		ID:       mainChannel.ID,
		Messages: []models.RoomChannelMessage{},
	}); err != nil {
		h.Collections.RoomCollection.DeleteOne(r.Context(), bson.M{"_id": inserted.InsertedID.(primitive.ObjectID)})
		responseMessage(w, http.StatusInternalServerError, "Internal error")
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inserted.InsertedID.(primitive.ObjectID).Hex())
}

func (h handler) UpdateRoom(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetUserFromRequest(r, r.Context(), *h.Collections, h.RedisClient)
	if err != nil {
		responseMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		responseMessage(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	var roomInput validation.Room
	if json.Unmarshal(body, &roomInput); err != nil {
		responseMessage(w, http.StatusBadRequest, "Bad request")
		return
	}
	validate := validator.New()
	if err := validate.Struct(roomInput); err != nil {
		responseMessage(w, http.StatusBadRequest, "Bad request")
		return
	}

	var room models.Room
	if err := h.Collections.RoomCollection.FindOne(r.Context(), bson.M{"_id": id}).Decode(&room); err != nil {
		if err != mongo.ErrNoDocuments {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
			return
		}
	}
	if room.Author != user.ID {
		responseMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if _, err := h.Collections.RoomCollection.UpdateByID(r.Context(), id, bson.M{
		"$set": bson.M{
			"name":       roomInput.Name,
			"is_private": roomInput.Private,
		},
	}); err != nil {
		responseMessage(w, http.StatusInternalServerError, "Internal error")
	} else {
		responseMessage(w, http.StatusOK, "Room updated")
	}
}

// One API route for updating/inserting/deleting room channel data... maybe not best practice?
func (h handler) UpdateRoomChannelsData(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetUserFromRequest(r, r.Context(), *h.Collections, h.RedisClient)
	if err != nil {
		responseMessage(w, http.StatusUnauthorized, "StatusUnauthorized")
		return
	}

	roomId, err := primitive.ObjectIDFromHex(mux.Vars(r)["roomId"])
	if err != nil {
		responseMessage(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var roomInternalData models.RoomInternalData
	if err := h.Collections.RoomInternalDataCollection.FindOne(r.Context(), bson.M{"_id": roomId}).Decode(&roomInternalData); err != nil {
		if err != mongo.ErrNoDocuments {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
		} else {
			responseMessage(w, http.StatusNotFound, "Not found")
		}
		return
	}

	var room models.Room
	if err := h.Collections.RoomCollection.FindOne(r.Context(), bson.M{"_id": roomId}).Decode(&room); err != nil {
		if err != mongo.ErrNoDocuments {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
			return
		}
	}
	if room.Author != user.ID {
		responseMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseMessage(w, http.StatusBadRequest, "Bad request")
		return
	}
	var updateRoomChannelsData validation.UpdateRoomChannelsData
	if err := json.Unmarshal(body, &updateRoomChannelsData); err != nil {
		responseMessage(w, http.StatusBadRequest, "Bad request")
		return
	}
	validate := validator.New()
	if err := validate.Struct(updateRoomChannelsData); err != nil {
		responseMessage(w, http.StatusBadRequest, "Bad request")
		return
	}

	// Update room channel names
	for _, urcd := range updateRoomChannelsData.UpdateData {
		id, err := primitive.ObjectIDFromHex(urcd.ID)
		if err != nil {
			responseMessage(w, http.StatusBadRequest, "Invalid ID")
			return
		}
		for _, hex := range updateRoomChannelsData.Delete {
			oid, err := primitive.ObjectIDFromHex(hex)
			if err != nil {
				responseMessage(w, http.StatusBadRequest, "Invalid ID")
				return
			}
			// Don't bother updating channels that are going to be deleted in the next step
			if oid == id {
				continue
			}
		}
		res, err := h.Collections.RoomChannelCollection.UpdateOne(r.Context(), bson.M{"_id": id, "room_id": roomId}, bson.M{
			"$set": bson.M{
				"name": strings.TrimSpace(urcd.Name),
			},
		})
		if err != nil {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
			return
		}
		if res.MatchedCount == 0 {
			responseMessage(w, http.StatusBadRequest, "Bad request")
			return
		}
	}

	insertedNewMainChannel := false
	insertedChannels := []models.RoomChannel{}

	// Insert room channels. The user can create multiple channels with the same name, but it doesn't really matter
	for _, insertData := range updateRoomChannelsData.InsertData {
		res, err := h.Collections.RoomChannelCollection.InsertOne(r.Context(), models.RoomChannel{
			ID:     primitive.NewObjectID(),
			RoomID: roomId,
			Name:   strings.TrimSpace(insertData.Name),
		})
		if err != nil {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
			return
		} else {
			if _, err := h.Collections.RoomChannelMessagesCollection.InsertOne(r.Context(), models.RoomChannelMessages{
				ID:       res.InsertedID.(primitive.ObjectID),
				Messages: []models.RoomChannelMessage{},
			}); err != nil {
				responseMessage(w, http.StatusInternalServerError, "Internal error")
				return
			}
		}
		if _, err := h.Collections.RoomInternalDataCollection.UpdateByID(r.Context(), roomId, bson.M{
			"$push": bson.M{
				"channels": res.InsertedID.(primitive.ObjectID),
			},
		}); err != nil {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
			return
		}
		if insertData.PromoteToMain {
			insertedNewMainChannel = true
		}
		if insertData.PromoteToMain {
			if _, err := h.Collections.RoomInternalDataCollection.UpdateByID(r.Context(), roomId, bson.M{
				"$set": bson.M{
					"main_channel": res.InsertedID.(primitive.ObjectID),
				},
			}); err != nil {
				responseMessage(w, http.StatusInternalServerError, "Internal error")
				return
			}
		}
		insertedChannels = append(insertedChannels, models.RoomChannel{
			ID:     res.InsertedID.(primitive.ObjectID),
			RoomID: roomId,
			Name:   strings.TrimSpace(insertData.Name),
		})
	}

	// Delete room channels
	if int64(len(updateRoomChannelsData.Delete)) != 0 {
		ids := []primitive.ObjectID{}
		for _, v := range updateRoomChannelsData.Delete {
			id, err := primitive.ObjectIDFromHex(v)
			if err != nil {
				responseMessage(w, http.StatusBadRequest, "Invalid ID")
				return
			}
			ids = append(ids, id)
		}
		_, err := h.Collections.RoomChannelCollection.UpdateMany(r.Context(), bson.M{
			"_id": bson.M{
				"$in": ids,
			},
			"room_id": roomId,
		}, bson.M{"$set": bson.M{"to_be_deleted": true}})
		if err != nil {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
			return
		}
	}

	// Promote pre-existing room channel to main (step will be skipped if a new main channel was already created)
	if !insertedNewMainChannel && updateRoomChannelsData.PromoteToMain != "" {
		promoteToMain, err := primitive.ObjectIDFromHex(updateRoomChannelsData.PromoteToMain)
		if err == nil {
			if promoteToMain != primitive.NilObjectID {
				found := false
				oid := &primitive.ObjectID{}
				for _, oi := range roomInternalData.Channels {
					if oi == promoteToMain {
						oid = &oi
						found = true
						break
					}
				}
				if found {
					if _, err := h.Collections.RoomInternalDataCollection.UpdateByID(r.Context(), roomId, bson.M{
						"$set": bson.M{
							"main_channel": oid,
						},
					}); err != nil {
						responseMessage(w, http.StatusInternalServerError, "Internal error")
						return
					}
				}
			}
		}
	}

	responseMessage(w, http.StatusOK, "Room channels data updated")
}

func (h handler) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetUserFromRequest(r, r.Context(), *h.Collections, h.RedisClient)
	if err != nil {
		responseMessage(w, http.StatusUnauthorized, "StatusUnauthorized")
		return
	}

	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		responseMessage(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var room models.Room
	if err := h.Collections.RoomCollection.FindOne(r.Context(), bson.M{"_id": id}).Decode(&room); err != nil {
		if err != mongo.ErrNoDocuments {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
			return
		}
	}
	if room.Author != user.ID {
		responseMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if _, err := h.Collections.RoomCollection.DeleteOne(r.Context(), bson.M{"_id": id}); err != nil {
		responseMessage(w, http.StatusInternalServerError, "Internal error")
		return
	}

	responseMessage(w, http.StatusOK, "Room deleted")
}

func (h handler) UploadRoomImage(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		responseMessage(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	user, err := helpers.GetUserFromRequest(r, r.Context(), *h.Collections, h.RedisClient)
	if err != nil {
		responseMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	room := &models.Room{}
	if err := h.Collections.RoomCollection.FindOne(r.Context(), bson.M{"_id": id}).Decode(&room); err != nil {
		if err == mongo.ErrNoDocuments {
			responseMessage(w, http.StatusNotFound, "Not found")
		} else {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
		}
		return
	}

	if room.Author != user.ID {
		responseMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	r.ParseMultipartForm(32 << 20) // binary shift maximum 20mb in bytes
	file, handler, err := r.FormFile("file")
	defer file.Close()
	if err != nil {
		responseMessage(w, http.StatusBadRequest, "Bad request")
		return
	}

	if handler.Size > 20*1024*1024 {
		responseMessage(w, http.StatusRequestEntityTooLarge, "File too large. Max 20mb")
		return
	}

	src, err := handler.Open()
	if err != nil {
		responseMessage(w, http.StatusBadRequest, "Bad request")
	}

	var isJPEG, isPNG bool
	isJPEG = handler.Header.Get("Content-Type") == "image/jpeg"
	isPNG = handler.Header.Get("Content-Type") == "image/png"
	if !isJPEG && !isPNG {
		responseMessage(w, http.StatusBadRequest, "Unrecognized file format")
		return
	}

	var img image.Image
	var blurImg image.Image
	var decodeErr error
	if isJPEG {
		img, decodeErr = jpeg.Decode(src)
	} else {
		img, decodeErr = png.Decode(src)
	}
	if decodeErr != nil {
		responseMessage(w, http.StatusInternalServerError, "Internal error")
		return
	}
	blurImg = img
	buf := &bytes.Buffer{}
	blurBuf := &bytes.Buffer{}
	width := img.Bounds().Dx()
	if width > 350 {
		img = resize.Resize(350, 0, img, resize.Lanczos2)
	} else {
		img = resize.Resize(uint(width), 0, img, resize.Lanczos2)
	}
	blurImg = resize.Resize(10, 0, img, resize.Bilinear)
	if err := jpeg.Encode(buf, img, nil); err != nil {
		responseMessage(w, http.StatusInternalServerError, "Internal error")
		return
	}
	if err := jpeg.Encode(blurBuf, blurImg, nil); err != nil {
		responseMessage(w, http.StatusInternalServerError, "Internal error")
		return
	}

	if _, err := h.Collections.RoomCollection.UpdateByID(r.Context(), id, bson.M{
		"$set": bson.M{
			"blur": "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(blurBuf.Bytes()),
		},
	}); err != nil {
		responseMessage(w, http.StatusInternalServerError, "Internal error")
		return
	}

	foundImage := false
	res := h.Collections.RoomImageCollection.FindOne(r.Context(), bson.M{"_id": id})
	if res.Err() != nil {
		if res.Err() != mongo.ErrNoDocuments {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
			return
		}
	} else {
		foundImage = true
	}

	if foundImage {
		if _, err := h.Collections.RoomImageCollection.UpdateByID(r.Context(), id, bson.M{
			"$set": bson.M{
				"binary": buf.Bytes(),
			},
		}); err != nil {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
			return
		}
	} else {
		if _, err := h.Collections.RoomImageCollection.InsertOne(r.Context(), models.RoomImage{
			ID:     id,
			Binary: primitive.Binary{Data: buf.Bytes()},
		}); err != nil {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
			return
		}
	}

	if foundImage {
		responseMessage(w, http.StatusOK, "Image updated")
	} else {
		responseMessage(w, http.StatusCreated, "Image created")
	}

	outBytes, err := json.Marshal(socketmodels.OutChangeMessage{
		Type:   "CHANGE",
		Method: "UPDATE_IMAGE",
		Entity: "ROOM",
		Data:   `{"ID":"` + room.ID.Hex() + `"}`,
	})
	if err != nil {
		responseMessage(w, http.StatusInternalServerError, "Internal error")
		return
	}

	h.SocketServer.SendDataToSubscription <- socketserver.SubscriptionDataMessage{
		Name: "room-display-data=" + room.ID.Hex(),
		Data: outBytes,
	}
}

func (h handler) GetRoomPage(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetUserFromRequest(r, r.Context(), *h.Collections, h.RedisClient)
	if err != nil {
		responseMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	filter := bson.M{}
	if r.URL.Query().Has("own") {
		filter = bson.M{"author": user.ID}
	}

	pageNumberString := mux.Vars(r)["page"]
	pageNumber, err := strconv.Atoi(pageNumberString)
	if err != nil {
		responseMessage(w, http.StatusBadRequest, "Invalid page")
		return
	}
	pageSize := 20

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSkip(int64(pageSize) * (int64(pageNumber) - 1))

	rooms := []models.Room{}
	if cursor, err := h.Collections.RoomCollection.Find(r.Context(), filter, findOptions); err != nil {
		cursor.Close(r.Context())
		if err == mongo.ErrNoDocuments {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(rooms)
		} else {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
		}
		return
	} else {
		for cursor.Next(r.Context()) {
			room := &models.Room{}
			if err := cursor.Decode(&room); err != nil {
				cursor.Close(r.Context())
				responseMessage(w, http.StatusInternalServerError, "Internal error")
				return
			}
			rooms = append(rooms, *room)
		}
		cursor.Close(r.Context())
	}

	outRooms := []models.Room{}
	for _, room := range rooms {
		var roomExternalData models.RoomExternalData
		if err := h.Collections.RoomExternalDataCollection.FindOne(r.Context(), bson.M{"_id": room.ID}).Decode(&roomExternalData); err != nil {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
			return
		}
		addRoom := true
		// If user is banned from a room don't add it to the response
		for _, oi := range roomExternalData.Banned {
			if oi == user.ID {
				addRoom = false
				break
			}
		}
		isMember := room.Author == user.ID
		if !isMember {
			// If user is not a member, and the room is private, don't add it to the response
			for _, oi := range roomExternalData.Members {
				if oi == user.ID {
					isMember = true
					break
				}
			}
		}
		if !addRoom || roomExternalData.Private && !isMember {
			continue
		}
		outRooms = append(outRooms, room)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(outRooms)
}

func (h handler) GetOwnRoomIDs(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetUserFromRequest(r, r.Context(), *h.Collections, h.RedisClient)
	if err != nil {
		responseMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	cursor, err := h.Collections.RoomCollection.Find(r.Context(), bson.M{"author": user.ID})
	defer cursor.Close(r.Context())
	rooms := []primitive.ObjectID{}
	if err != nil {
		responseMessage(w, http.StatusInternalServerError, "Internal error")
		return
	}
	for cursor.Next(r.Context()) {
		room := &models.Room{}
		cursor.Decode(&room)
		rooms = append(rooms, room.ID)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rooms)
}

func (h handler) GetRoom(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetUserFromRequest(r, r.Context(), *h.Collections, h.RedisClient)
	if err != nil {
		responseMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		responseMessage(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	room := &models.Room{}
	if err := h.Collections.RoomCollection.FindOne(r.Context(), bson.M{"_id": id}).Decode(&room); err != nil {
		if err != mongo.ErrNoDocuments {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
		} else {
			responseMessage(w, http.StatusNotFound, "Room not found")
		}
		return
	}

	roomExternalData := &models.RoomExternalData{}
	if err := h.Collections.RoomExternalDataCollection.FindOne(r.Context(), bson.M{"_id": id}).Decode(&roomExternalData); err != nil {
		if err != mongo.ErrNoDocuments {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
		} else {
			responseMessage(w, http.StatusNotFound, "Room not found")
		}
		return
	}

	roomInternalData := &models.RoomInternalData{}
	if err := h.Collections.RoomInternalDataCollection.FindOne(r.Context(), bson.M{"_id": id}).Decode(&roomInternalData); err != nil {
		if err != mongo.ErrNoDocuments {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
		} else {
			responseMessage(w, http.StatusNotFound, "Room not found")
		}
		return
	}

	if room.Author != user.ID {
		for _, oi := range roomExternalData.Banned {
			if oi == user.ID {
				responseMessage(w, http.StatusUnauthorized, "You are banned from this room")
				return
			}
		}
		if roomExternalData.Private {
			isMember := false
			for _, oi := range roomExternalData.Members {
				if oi == user.ID {
					isMember = true
					break
				}
			}
			if !isMember {
				responseMessage(w, http.StatusForbidden, "You are not a member of this room")
				return
			}
		}
	}

	room.Private = roomExternalData.Private
	room.Members = roomExternalData.Members
	room.Banned = roomExternalData.Banned
	room.Channels = roomInternalData.Channels
	room.MainChannel = roomInternalData.MainChannel

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(room)
}

func (h handler) GetRoomDisplayData(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetUserFromRequest(r, r.Context(), *h.Collections, h.RedisClient)
	if err != nil {
		responseMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		responseMessage(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	room := &models.Room{}
	if err := h.Collections.RoomCollection.FindOne(r.Context(), bson.M{"_id": id}).Decode(&room); err != nil {
		if err != mongo.ErrNoDocuments {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
		} else {
			responseMessage(w, http.StatusNotFound, "Room not found")
		}
		return
	}

	roomExternalData := &models.RoomExternalData{}
	if err := h.Collections.RoomExternalDataCollection.FindOne(r.Context(), bson.M{"_id": id}).Decode(&roomExternalData); err != nil {
		if err != mongo.ErrNoDocuments {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
		} else {
			responseMessage(w, http.StatusNotFound, "Room not found")
		}
		return
	}

	if room.Author != user.ID {
		for _, oi := range roomExternalData.Banned {
			if oi == user.ID {
				responseMessage(w, http.StatusUnauthorized, "You are banned from this room")
				return
			}
		}
		if roomExternalData.Private {
			isMember := false
			for _, oi := range roomExternalData.Members {
				if oi == user.ID {
					isMember = true
					break
				}
			}
			if !isMember {
				responseMessage(w, http.StatusForbidden, "You are not a member of this room")
				return
			}
		}
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(room)
}

func (h handler) GetRoomImage(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetUserFromRequest(r, r.Context(), *h.Collections, h.RedisClient)
	if err != nil {
		responseMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		responseMessage(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	roomExternalData := &models.RoomExternalData{}
	if err := h.Collections.RoomExternalDataCollection.FindOne(r.Context(), bson.M{"_id": id}).Decode(&roomExternalData); err != nil {
		if err == mongo.ErrNoDocuments {
			responseMessage(w, http.StatusNotFound, "Not found")
		} else {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
		}
		return
	}
	room := &models.Room{}
	if err := h.Collections.RoomCollection.FindOne(r.Context(), bson.M{"_id": id}).Decode(&room); err != nil {
		if err == mongo.ErrNoDocuments {
			responseMessage(w, http.StatusNotFound, "Not found")
		} else {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
		}
		return
	}

	for _, oi := range roomExternalData.Banned {
		if oi == user.ID {
			responseMessage(w, http.StatusUnauthorized, "You are banned from this room")
			return
		}
	}

	if roomExternalData.Private && room.Author != user.ID {
		isMember := false
		for _, oi := range roomExternalData.Members {
			if oi == user.ID {
				isMember = true
				break
			}
		}
		if !isMember {
			responseMessage(w, http.StatusForbidden, "You are not a member of this room")
			return
		}
	}

	roomImage := &models.RoomImage{}
	if err := h.Collections.RoomImageCollection.FindOne(r.Context(), bson.M{"_id": id}).Decode(&roomImage); err != nil {
		if err == mongo.ErrNoDocuments {
			responseMessage(w, http.StatusNotFound, "Not found")
		} else {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
		}
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(roomImage.Binary.Data)))
	if _, err := w.Write(roomImage.Binary.Data); err != nil {
		log.Println("Unable to write image to response")
	}
}

func (h handler) GetRoomChannelsData(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetUserFromRequest(r, r.Context(), *h.Collections, h.RedisClient)
	if err != nil {
		responseMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	roomId, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		responseMessage(w, http.StatusBadRequest, "Invalid room ID")
		return
	}

	roomExternalData := &models.RoomExternalData{}
	if err := h.Collections.RoomExternalDataCollection.FindOne(r.Context(), bson.M{"_id": roomId}).Decode(&roomExternalData); err != nil {
		if err == mongo.ErrNoDocuments {
			responseMessage(w, http.StatusNotFound, "Not found")
		} else {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
		}
		return
	}
	roomInternalData := &models.RoomInternalData{}
	if err := h.Collections.RoomInternalDataCollection.FindOne(r.Context(), bson.M{"_id": roomId}).Decode(&roomInternalData); err != nil {
		if err == mongo.ErrNoDocuments {
			responseMessage(w, http.StatusNotFound, "Not found")
		} else {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
		}
		return
	}
	room := &models.Room{}
	if err := h.Collections.RoomCollection.FindOne(r.Context(), bson.M{"_id": roomId}).Decode(&room); err != nil {
		if err == mongo.ErrNoDocuments {
			responseMessage(w, http.StatusNotFound, "Not found")
		} else {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
		}
		return
	}

	for _, oi := range roomExternalData.Banned {
		if oi == user.ID {
			responseMessage(w, http.StatusUnauthorized, "You are banned from this room")
			return
		}
	}

	if roomExternalData.Private && room.Author != user.ID {
		isMember := false
		for _, oi := range roomExternalData.Members {
			if oi == user.ID {
				isMember = true
				break
			}
		}
		if !isMember {
			responseMessage(w, http.StatusForbidden, "You are not a member of this room")
			return
		}
	}

	channels := []models.RoomChannel{}
	if cursor, err := h.Collections.RoomChannelCollection.Find(r.Context(), bson.M{
		"_id": bson.M{
			"$in": roomInternalData.Channels,
		},
		"room_id": roomId,
	}); err != nil {
		log.Println(err)
		cursor.Close(r.Context())
		if err == mongo.ErrNoDocuments {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(channels)
		} else {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
		}
		return
	} else {
		for cursor.Next(r.Context()) {
			channel := &models.RoomChannel{}
			err := cursor.Decode(&channel)
			if err != nil {
				cursor.Close(r.Context())
				responseMessage(w, http.StatusInternalServerError, "Internal error")
				return
			}
			channels = append(channels, *channel)
		}
		cursor.Close(r.Context())
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(channels)
}

func (h handler) GetRoomChannel(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetUserFromRequest(r, r.Context(), *h.Collections, h.RedisClient)
	if err != nil {
		responseMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		responseMessage(w, http.StatusBadRequest, "Bad request")
		return
	}
	roomId, err := primitive.ObjectIDFromHex(mux.Vars(r)["roomId"])
	if err != nil {
		responseMessage(w, http.StatusBadRequest, "Bad request")
		return
	}

	room := &models.Room{}
	if err := h.Collections.RoomCollection.FindOne(r.Context(), bson.M{"_id": roomId}).Decode(&room); err != nil {
		if err == mongo.ErrNoDocuments {
			responseMessage(w, http.StatusNotFound, "Not found")
		} else {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
		}
		return
	}

	roomExternalData := &models.RoomExternalData{}
	if err := h.Collections.RoomExternalDataCollection.FindOne(r.Context(), bson.M{"_id": roomId}).Decode(&roomExternalData); err != nil {
		if err == mongo.ErrNoDocuments {
			responseMessage(w, http.StatusNotFound, "Room not found")
		} else {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
		}
		return
	}

	for _, oi := range roomExternalData.Banned {
		if oi == user.ID {
			responseMessage(w, http.StatusUnauthorized, "You are banned from this room")
			return
		}
	}

	if roomExternalData.Private && room.Author != user.ID {
		isMember := false
		for _, oi := range roomExternalData.Members {
			if oi == user.ID {
				isMember = true
				break
			}
		}
		if !isMember {
			responseMessage(w, http.StatusForbidden, "You are not a member of this room")
			return
		}
	}

	roomChannel := &models.RoomChannel{}
	if err := h.Collections.RoomChannelCollection.FindOne(r.Context(), bson.M{"_id": id}).Decode(&roomChannel); err != nil {
		if err == mongo.ErrNoDocuments {
			responseMessage(w, http.StatusNotFound, "Channel not found")
		} else {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
		}
		return
	}

	if roomChannel.RoomID != roomId {
		responseMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	roomChannelMessages := &models.RoomChannelMessages{}
	if err := h.Collections.RoomChannelMessagesCollection.FindOne(r.Context(), bson.M{"_id": id}).Decode(&roomChannelMessages); err != nil {
		if err == mongo.ErrNoDocuments {
			responseMessage(w, http.StatusNotFound, "Channel not found")
		} else {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
		}
		return
	}

	roomChannel.Messages = roomChannelMessages.Messages

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(roomChannel)
}
