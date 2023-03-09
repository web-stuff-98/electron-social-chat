package handlers

/* Dependency injection for handlers */

import (
	"encoding/json"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/web-stuff-98/electron-social-chat/pkg/attachmentserver"
	"github.com/web-stuff-98/electron-social-chat/pkg/db"
	"github.com/web-stuff-98/electron-social-chat/pkg/socketserver"

	"go.mongodb.org/mongo-driver/mongo"
)

func responseMessage(w http.ResponseWriter, c int, m string) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(c)
	json.NewEncoder(w).Encode(map[string]string{"msg": m})
	return
}

type handler struct {
	DB               *mongo.Database
	Collections      *db.Collections
	RedisClient      *redis.Client
	SocketServer     *socketserver.SocketServer
	AttachmentServer *attachmentserver.AttachmentServer
}

func New(db *mongo.Database, collections *db.Collections, redisClient *redis.Client, socketServer *socketserver.SocketServer, attachmentServer *attachmentserver.AttachmentServer) handler {
	return handler{db, collections, redisClient, socketServer, attachmentServer}
}
