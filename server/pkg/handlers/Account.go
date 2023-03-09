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

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
	"github.com/web-stuff-98/electron-social-chat/pkg/db/models"
	"github.com/web-stuff-98/electron-social-chat/pkg/helpers"
	"github.com/web-stuff-98/electron-social-chat/pkg/validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func (h handler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseMessage(w, http.StatusBadRequest, "Bad request")
		return
	}
	var credentialsInput validation.Credentials
	if json.Unmarshal(body, &credentialsInput); err != nil {
		responseMessage(w, http.StatusBadRequest, "Bad request")
		return
	}
	validate := validator.New()
	if err := validate.Struct(credentialsInput); err != nil {
		responseMessage(w, http.StatusBadRequest, "Bad request")
		return
	}

	var user models.User
	if err := h.Collections.UserCollection.FindOne(r.Context(), bson.M{
		"username": bson.M{
			"$regex":   credentialsInput.Username,
			"$options": "i",
		},
	}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			responseMessage(w, http.StatusNotFound, "No account exists by that name")
		} else {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentialsInput.Password)); err != nil {
		responseMessage(w, http.StatusUnauthorized, "Incorrect credentials")
		return
	}

	if cookie, err := helpers.GenerateCookieAndSession(r.Context(), user.ID, *h.Collections, h.RedisClient); err != nil {
		responseMessage(w, http.StatusBadRequest, err.Error())
		return
	} else {
		http.SetCookie(w, &cookie)

		user.IsOnline = true

		var pfp models.Pfp
		if err := h.Collections.PfpCollection.FindOne(r.Context(), bson.M{"_id": user.ID}).Decode(&pfp); err != nil {
			if err != mongo.ErrNoDocuments {
				responseMessage(w, http.StatusInternalServerError, "Internal error")
				return
			}
		} else {
			user.Base64pfp = "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(pfp.Binary.Data)
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}

func (h handler) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseMessage(w, http.StatusBadRequest, "Bad request")
		return
	}
	var credentialsInput validation.Credentials
	if err := json.Unmarshal(body, &credentialsInput); err != nil {
		responseMessage(w, http.StatusBadRequest, "Bad request")
		return
	}
	validate := validator.New()
	if err := validate.Struct(credentialsInput); err != nil {
		responseMessage(w, http.StatusBadRequest, "Bad request")
		return
	}

	res := h.Collections.UserCollection.FindOne(r.Context(), bson.M{
		"username": bson.M{
			"$regex":   credentialsInput.Username,
			"$options": "i",
		},
	})
	if res.Err() != nil {
		if res.Err() != mongo.ErrNoDocuments {
			responseMessage(w, http.StatusBadRequest, "Internal error")
			return
		}
	} else {
		responseMessage(w, http.StatusBadRequest, "There is another user by that name already")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(credentialsInput.Password), 14)
	if err != nil {
		responseMessage(w, http.StatusInternalServerError, "Internal error")
		return
	}

	user := models.User{
		ID:       primitive.NewObjectID(),
		Username: credentialsInput.Username,
		Password: string(hash),
		IsOnline: true,
	}

	userMessagingData := models.UserMessagingData{
		ID:                   user.ID,
		Messages:             []models.DirectMessage{},
		MessagesSentTo:       []primitive.ObjectID{},
		MessagesReceivedFrom: []primitive.ObjectID{},
		Blocked:              []primitive.ObjectID{},
		Friends:              []primitive.ObjectID{},
		FriendRequests:       []models.FriendRequest{},
		Invitations:          []models.Invitation{},
	}

	if _, err := h.Collections.UserCollection.InsertOne(r.Context(), user); err != nil {
		responseMessage(w, http.StatusInternalServerError, "Internal error")
		return
	}

	if _, err := h.Collections.UserMessagingDataCollection.InsertOne(r.Context(), userMessagingData); err != nil {
		responseMessage(w, http.StatusInternalServerError, "Internal error")
		return
	}

	if cookie, err := helpers.GenerateCookieAndSession(r.Context(), user.ID, *h.Collections, h.RedisClient); err != nil {
		responseMessage(w, http.StatusBadRequest, err.Error())
		return
	} else {
		http.SetCookie(w, &cookie)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}

func (h handler) Refresh(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetUserFromRequest(r, r.Context(), *h.Collections, h.RedisClient)
	if err != nil {
		responseMessage(w, http.StatusForbidden, "Forbidden")
		clearedCookie := helpers.GetClearedCookie()
		http.SetCookie(w, &clearedCookie)
		return
	}

	if cookie, err := helpers.GenerateCookieAndSession(r.Context(), user.ID, *h.Collections, h.RedisClient); err != nil {
		responseMessage(w, http.StatusBadRequest, err.Error())
		clearedCookie := helpers.GetClearedCookie()
		http.SetCookie(w, &clearedCookie)
		return
	} else {
		http.SetCookie(w, &cookie)
		responseMessage(w, http.StatusOK, "Token refreshed")
	}
}

func (h handler) Logout(w http.ResponseWriter, r *http.Request) {
	clearedCookie := helpers.GetClearedCookie()
	http.SetCookie(w, &clearedCookie)
	responseMessage(w, http.StatusOK, "Logged out")
}

func (h handler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetUserFromRequest(r, r.Context(), *h.Collections, h.RedisClient)
	if err != nil {
		responseMessage(w, http.StatusForbidden, "Forbidden")
		clearedCookie := helpers.GetClearedCookie()
		http.SetCookie(w, &clearedCookie)
		return
	}

	if res, err := h.Collections.UserCollection.DeleteOne(r.Context(), bson.M{"_id": user.ID}); err != nil {
		responseMessage(w, http.StatusInternalServerError, "Internal error")
		return
	} else {
		clearedCookie := helpers.GetClearedCookie()
		http.SetCookie(w, &clearedCookie)
		if res.DeletedCount == 0 {
			responseMessage(w, http.StatusNotFound, "Your account does not exist")
			return
		}
		responseMessage(w, http.StatusOK, "Account deleted")
		return
	}
}

func (h handler) UploadPfp(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetUserFromRequest(r, r.Context(), *h.Collections, h.RedisClient)
	if err != nil {
		responseMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	r.ParseMultipartForm(32 << 20) // what is <<, binary shift whatever that is. Is used here to define max size in bytes (20mb)

	file, handler, err := r.FormFile("file")
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
	}
	defer file.Close()

	if handler.Size > 20*1024*1024 {
		responseMessage(w, http.StatusRequestEntityTooLarge, "File too large, max 20mb.")
		return
	}

	src, err := handler.Open()
	if err != nil {
		responseMessage(w, http.StatusBadRequest, "Bad request")
		return
	}
	var isJPEG, isPNG bool
	isJPEG = handler.Header.Get("Content-Type") == "image/jpeg"
	isPNG = handler.Header.Get("Content-Type") == "image/png"
	if !isJPEG && !isPNG {
		responseMessage(w, http.StatusBadRequest, "Only JPEG and PNG are supported")
		return
	}
	var img image.Image
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
	buf := &bytes.Buffer{}
	if img.Bounds().Dx() > img.Bounds().Dy() {
		img = resize.Resize(64, 0, img, resize.Lanczos3)
	} else {
		img = resize.Resize(0, 64, img, resize.Lanczos3)
	}
	if err := jpeg.Encode(buf, img, nil); err != nil {
		responseMessage(w, http.StatusInternalServerError, "Internal error")
		return
	}
	pfpBytes := buf.Bytes()
	base64pfp := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(pfpBytes)
	user.Base64pfp = base64pfp

	res, err := h.Collections.PfpCollection.UpdateByID(r.Context(), user.ID, bson.M{"$set": bson.M{"binary": primitive.Binary{Data: buf.Bytes()}}})
	if err != nil {
		responseMessage(w, http.StatusInternalServerError, "Internal error")
		return
	}
	if res.MatchedCount == 0 {
		_, err := h.Collections.PfpCollection.InsertOne(r.Context(), models.Pfp{
			ID:     user.ID,
			Binary: primitive.Binary{Data: buf.Bytes()},
		})
		if err != nil {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
			return
		}
	}

	responseMessage(w, http.StatusOK, "Pfp updated")
}

func (h handler) GetConversation(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetUserFromRequest(r, r.Context(), *h.Collections, h.RedisClient)
	if err != nil {
		responseMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	converseId, err := primitive.ObjectIDFromHex(mux.Vars(r)["uid"])
	if err != nil {
		responseMessage(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	messagingData := &models.UserMessagingData{}
	if h.Collections.UserMessagingDataCollection.FindOne(r.Context(), bson.M{"_id": user.ID}).Decode(&messagingData); err != nil {
		if err == mongo.ErrNoDocuments {
			responseMessage(w, http.StatusNotFound, "Not found")
		} else {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
		}
		return
	}

	converseMessagingData := &models.UserMessagingData{}
	if h.Collections.UserMessagingDataCollection.FindOne(r.Context(), bson.M{"_id": converseId}).Decode(&converseMessagingData); err != nil {
		if err == mongo.ErrNoDocuments {
			responseMessage(w, http.StatusNotFound, "Not found")
		} else {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
		}
		return
	}

	messages := []models.DirectMessage{}
	for _, dm := range messagingData.Messages {
		if dm.Author == converseId {
			messages = append(messages, dm)
		}
	}
	for _, dm := range converseMessagingData.Messages {
		if dm.Author == user.ID {
			messages = append(messages, dm)
		}
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
}

func (h handler) GetConversations(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetUserFromRequest(r, r.Context(), *h.Collections, h.RedisClient)
	if err != nil {
		responseMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	messagingData := &models.UserMessagingData{}
	if h.Collections.UserMessagingDataCollection.FindOne(r.Context(), bson.M{"_id": user.ID}).Decode(&messagingData); err != nil {
		if err == mongo.ErrNoDocuments {
			responseMessage(w, http.StatusNotFound, "Not found")
		} else {
			responseMessage(w, http.StatusInternalServerError, "Internal error")
		}
		return
	}

	conversationsUnique := make(map[primitive.ObjectID]struct{})
	conversations := []primitive.ObjectID{}
	for _, oi := range messagingData.MessagesReceivedFrom {
		conversationsUnique[oi] = struct{}{}
	}
	for _, oi := range messagingData.MessagesSentTo {
		conversationsUnique[oi] = struct{}{}
	}
	for oi := range conversationsUnique {
		conversations = append(conversations, oi)
	}

	invitations := []models.Invitation{}
	friendRequests := []models.FriendRequest{}
	invitations = append(invitations, messagingData.Invitations...)
	friendRequests = append(friendRequests, messagingData.FriendRequests...)

	sentAndReceived := []primitive.ObjectID{}
	for _, oi := range messagingData.MessagesSentTo {
		sentAndReceived = append(sentAndReceived, oi)
	}
	for _, oi := range messagingData.MessagesReceivedFrom {
		sentAndReceived = append(sentAndReceived, oi)
	}

	if cursor, err := h.Collections.UserMessagingDataCollection.Find(r.Context(), bson.M{"id": bson.M{"$in": sentAndReceived}}); err != nil {
		cursor.Close(r.Context())
		responseMessage(w, http.StatusInternalServerError, "Internal error")
		return
	} else {
		otherUsersMessagingData := []models.UserMessagingData{}
		if err := cursor.All(r.Context(), &otherUsersMessagingData); err != nil {
			cursor.Close(r.Context())
			responseMessage(w, http.StatusInternalServerError, "Internal error")
			return
		}
		log.Println("LEN:", len(otherUsersMessagingData))
		for _, umd := range otherUsersMessagingData {
			for _, frq := range umd.FriendRequests {
				if frq.Author == user.ID {
					frq.Recipient = umd.ID
					friendRequests = append(friendRequests, frq)
				}
			}
			for _, inv := range umd.Invitations {
				if inv.Author == user.ID {
					inv.Recipient = umd.ID
					invitations = append(invitations, inv)
				}
			}
		}
		cursor.Close(r.Context())
	}

	out := make(map[string]interface{})
	out["conversations"] = conversations
	out["friend_requests"] = friendRequests
	out["invitations"] = invitations

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(out)
}
