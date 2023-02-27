package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/web-stuff-98/electron-social-chat/pkg/db/models"
	"github.com/web-stuff-98/electron-social-chat/pkg/helpers"
	"github.com/web-stuff-98/electron-social-chat/pkg/validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	mainChannel := models.RoomChannel{
		ID:   primitive.NewObjectID(),
		Name: "Main channel",
	}
	if _, err := h.Collections.RoomInternalDataCollection.InsertOne(r.Context(), models.RoomInternalData{
		ID:          inserted.InsertedID.(primitive.ObjectID),
		Channels:    []models.RoomChannel{mainChannel},
		MainChannel: mainChannel.ID,
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

	responseMessage(w, http.StatusCreated, "Room created")
}

func (h handler) UpdateRoom(w http.ResponseWriter, r *http.Request) {
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

	if _, err := h.Collections.RoomCollection.UpdateByID(r.Context(), bson.M{"_id": id}, bson.M{
		"name":    roomInput.Name,
		"private": roomInput.Private,
	}); err != nil {
		responseMessage(w, http.StatusInternalServerError, "Internal error")
	} else {
		responseMessage(w, http.StatusOK, "Room updated")
	}
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
