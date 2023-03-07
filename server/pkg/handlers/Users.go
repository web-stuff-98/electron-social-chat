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

func (h handler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	currentUser, err := helpers.GetUserFromRequest(r, r.Context(), *h.Collections, h.RedisClient)
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
	var searchInput validation.UserSearch
	if err := json.Unmarshal(body, &searchInput); err != nil {
		responseMessage(w, http.StatusBadRequest, "Bad request")
		return
	}
	validate := validator.New()
	if err := validate.Struct(searchInput); err != nil {
		responseMessage(w, http.StatusBadRequest, "Bad request")
		return
	}

	users := []string{}
	cursor, err := h.Collections.UserCollection.Find(r.Context(), bson.M{
		"$text": bson.M{
			"$search":        searchInput.Username,
			"$caseSensitive": false,
		},
	})
	defer cursor.Close(r.Context())
	for cursor.Next(r.Context()) {
		var user models.User
		cursor.Decode(&user)
		if currentUser.ID != user.ID {
			users = append(users, user.ID.Hex())
		}
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func (h handler) GetUser(w http.ResponseWriter, r *http.Request) {
	_, err := helpers.GetUserFromRequest(r, r.Context(), *h.Collections, h.RedisClient)
	if err != nil {
		responseMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	uid, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		responseMessage(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	user := &models.User{}
	if err := h.Collections.UserCollection.FindOne(r.Context(), bson.M{"_id": uid}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			responseMessage(w, http.StatusNotFound, "Not found")
		} else {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
		}
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
