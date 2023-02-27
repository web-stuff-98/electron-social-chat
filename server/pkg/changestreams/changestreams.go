package changestreams

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"

	"github.com/web-stuff-98/electron-social-chat/pkg/db/models"
	"github.com/web-stuff-98/electron-social-chat/pkg/socketmodels"
	"github.com/web-stuff-98/electron-social-chat/pkg/socketserver"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var deletePipeline = bson.D{
	{
		Key: "$match", Value: bson.D{
			{Key: "operationType", Value: "delete"},
		},
	},
}
var updatePipeline = bson.D{
	{
		Key: "$match", Value: bson.D{
			{Key: "operationType", Value: "update"},
		},
	},
}
var insertPipeline = bson.D{
	{
		Key: "$match", Value: bson.D{
			{Key: "operationType", Value: "insert"},
		},
	},
}

func WatchCollections(DB *mongo.Database, ss *socketserver.SocketServer) {
	go watchUserDeletes(DB, ss)
	go watchUserPfpUpdates(DB, ss)
}

func watchUserDeletes(db *mongo.Database, ss *socketserver.SocketServer) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic watching user deletes :", r)
		}
	}()
	cs, err := db.Collection("users").Watch(context.Background(), mongo.Pipeline{deletePipeline}, options.ChangeStream().SetFullDocument(options.UpdateLookup))
	if err != nil {
		log.Panicln("CS ERR : ", err.Error())
	}
	for cs.Next(context.Background()) {
		var changeEv struct {
			DocumentKey struct {
				ID primitive.ObjectID `bson:"_id"`
			} `bson:"documentKey"`
		}
		err := cs.Decode(&changeEv)
		if err != nil {
			log.Println("CS DECODE ERROR : ", err)
			return
		}
		uid := changeEv.DocumentKey.ID

		db.Collection("pfps").DeleteOne(context.Background(), bson.M{"_id": uid})
	}
}

func watchUserPfpUpdates(db *mongo.Database, ss *socketserver.SocketServer) {
	cs, err := db.Collection("pfps").Watch(context.Background(), mongo.Pipeline{updatePipeline}, options.ChangeStream().SetFullDocument(options.UpdateLookup))
	if err != nil {
		log.Panicln("CS ERR : ", err.Error())
	}
	for cs.Next(context.Background()) {
		var changeEv struct {
			DocumentKey struct {
				ID primitive.ObjectID `bson:"_id"`
			} `bson:"documentKey"`
			FullDocument models.Pfp `bson:"fullDocument"`
		}
		err := cs.Decode(&changeEv)
		if err != nil {
			log.Println("CS DECODE ERROR : ", err)
			return
		}
		uid := changeEv.DocumentKey.ID
		pfp := &changeEv.FullDocument
		if err != nil {
			log.Println("CS JSON MARSHAL ERROR : ", err)
			return
		}
		if err != nil {
			log.Println("CS JSON MARSHAL ERROR : ", err)
			return
		}
		pfpB64 := map[string]string{
			"ID":        uid.Hex(),
			"base64pfp": "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(pfp.Binary.Data),
		}
		jsonBytes, err := json.Marshal(pfpB64)
		if err != nil {
			log.Println("CS MARSHAL ERROR : ", err)
			return
		}

		outBytes, err := json.Marshal(socketmodels.OutChangeMessage{
			Type:   "CHANGE",
			Method: "UPDATE_IMAGE",
			Entity: "USER",
			Data:   string(jsonBytes),
		})

		ss.SendDataToSubscription <- socketserver.SubscriptionDataMessage{
			Name: "user=" + uid.Hex(),
			Data: outBytes,
		}
	}
}
