package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/web-stuff-98/electron-social-chat/pkg/attachmentserver"
	"github.com/web-stuff-98/electron-social-chat/pkg/changestreams"
	"github.com/web-stuff-98/electron-social-chat/pkg/db"
	"github.com/web-stuff-98/electron-social-chat/pkg/handlers"
	rdb "github.com/web-stuff-98/electron-social-chat/pkg/redis"
	"github.com/web-stuff-98/electron-social-chat/pkg/socketserver"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	DB, colls := db.Init()
	router := mux.NewRouter()
	redis := rdb.Init()
	socketServer, err := socketserver.Init(colls)
	attachmentServer := attachmentserver.Init(socketServer, colls)
	if err != nil {
		log.Fatal("Error setting up socket server: ", err)
	}

	h := handlers.New(DB, colls, redis, socketServer, attachmentServer)

	var origins []string
	if os.Getenv("PRODUCTION") == "true" {
		origins = []string{"https://electron-social-chat-backend.herokuapp.com"}
	} else {
		origins = []string{"http://localhost:1420", "http://localhost:8080"}
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PATCH", "DELETE"},
		AllowCredentials: true,
	})

	api := router.PathPrefix("/api/").Subrouter()

	api.HandleFunc("/acc/login", h.Login).Methods(http.MethodPost)
	api.HandleFunc("/acc/register", h.Register).Methods(http.MethodPost)
	api.HandleFunc("/acc/refresh", h.Refresh).Methods(http.MethodPost)
	api.HandleFunc("/acc/logout", h.Logout).Methods(http.MethodPost)
	api.HandleFunc("/acc/delete", h.DeleteAccount).Methods(http.MethodDelete)
	api.HandleFunc("/acc/pfp", h.UploadPfp).Methods(http.MethodPost)
	api.HandleFunc("/acc/conversation/{uid}", h.GetConversation).Methods(http.MethodGet)
	api.HandleFunc("/acc/conversations", h.GetConversations).Methods(http.MethodGet)

	api.HandleFunc("/user/search", h.SearchUsers).Methods(http.MethodPost)
	api.HandleFunc("/user/{id}", h.GetUser).Methods(http.MethodGet)

	api.HandleFunc("/room/create", h.CreateRoom).Methods(http.MethodPost)
	api.HandleFunc("/room/update/{id}", h.UpdateRoom).Methods(http.MethodPatch)
	api.HandleFunc("/room/channels/{id}", h.GetRoomChannelsData).Methods(http.MethodGet)
	api.HandleFunc("/room/channels/update/{roomId}", h.UpdateRoomChannelsData).Methods(http.MethodPatch)
	api.HandleFunc("/room/channel/{roomId}/{id}", h.GetRoomChannel).Methods(http.MethodGet)
	api.HandleFunc("/room/delete/{id}", h.DeleteRoom).Methods(http.MethodDelete)
	api.HandleFunc("/room/{id}", h.GetRoom).Methods(http.MethodGet)
	api.HandleFunc("/room/image/{id}", h.GetRoomImage).Methods(http.MethodGet)
	api.HandleFunc("/room/display/{id}", h.GetRoomDisplayData).Methods(http.MethodGet)
	api.HandleFunc("/room/image/{id}", h.UploadRoomImage).Methods(http.MethodPost)
	api.HandleFunc("/room/page/{page}", h.GetRoomPage).Methods(http.MethodGet)
	api.HandleFunc("/rooms/own/ids", h.GetOwnRoomIDs).Methods(http.MethodGet)

	api.HandleFunc("/attachment/chunk/{msgId}", h.UploadAttachmentChunk).Methods(http.MethodPost)
	api.HandleFunc("/attachment/meta/{msgId}", h.GetAttachmentMetadata).Methods(http.MethodGet)
	api.HandleFunc("/attachment/meta", h.CreateAttachmentMetadata).Methods(http.MethodPost)
	api.HandleFunc("/attachment/{msgId}", h.GetAttachment).Methods(http.MethodGet)

	api.HandleFunc("/ws", h.WebSocketEndpoint)

	log.Println("Watching collections...")
	changestreams.WatchCollections(DB, socketServer)

	log.Println("API open on port", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(fmt.Sprint(":", os.Getenv("PORT")), c.Handler(router)))
}
