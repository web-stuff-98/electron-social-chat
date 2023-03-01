package handlers

import (
	"bytes"
	"encoding/json"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
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
			"blur": blurBuf.Bytes(),
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
}
