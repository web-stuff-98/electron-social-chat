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
	go watchUserPfpUpdates(DB, ss)
	go watchRoomChannelUpdates(DB, ss)
	go watchRoomUpdates(DB, ss)

	go watchUserDeletes(DB, ss)
	go watchRoomDeletes(DB, ss)
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
		db.Collection("rooms").DeleteMany(context.Background(), bson.M{"author": uid})

		userMessagingData := &models.UserMessagingData{}
		if err := db.Collection("user_messaging_data").FindOne(context.Background(), bson.M{"_id": uid}).Decode(&userMessagingData); err == nil {
			for _, oi := range userMessagingData.MessagesSentTo {
				// for each inbox the user has sent a message to, remove their messages
				db.Collection("user_messaging_data").UpdateByID(context.Background(), oi, bson.M{
					"$pull": bson.M{
						"messages": bson.M{
							"author": uid,
						},
					},
				})
			}
			for _, oi := range userMessagingData.MessagesReceivedFrom {
				// for each inbox the user has received a message from, remove those messages
				db.Collection("user_messaging_data").UpdateByID(context.Background(), oi, bson.M{
					"$pull": bson.M{
						"messages": bson.M{
							"recipient": uid,
						},
					},
				})
			}
		}

		db.Collection("user_messaging_data").DeleteOne(context.Background(), bson.M{"_id": uid})

		ss.DestroySubscription <- "user=" + uid.Hex()
	}
}

func watchUserPfpUpdates(db *mongo.Database, ss *socketserver.SocketServer) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic watching pfp updates :", r)
		}
	}()
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
		if err != nil {
			continue
		}

		ss.SendDataToSubscription <- socketserver.SubscriptionDataMessage{
			Name: "user=" + uid.Hex(),
			Data: outBytes,
		}
	}
}

func watchRoomUpdates(db *mongo.Database, ss *socketserver.SocketServer) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic watching room updates :", r)
		}
	}()
	cs, err := db.Collection("rooms").Watch(context.Background(), mongo.Pipeline{updatePipeline}, options.ChangeStream().SetFullDocument(options.UpdateLookup))
	if err != nil {
		log.Panicln("CS ERR : ", err.Error())
	}
	for cs.Next(context.Background()) {
		var changeEv struct {
			DocumentKey struct {
				ID primitive.ObjectID `bson:"_id"`
			} `bson:"documentKey"`
			FullDocument models.Room `bson:"fullDocument"`
		}
		err := cs.Decode(&changeEv)
		if err != nil {
			log.Println("CS DECODE ERROR : ", err)
			return
		}

		jsonBytes, err := json.Marshal(changeEv.FullDocument)
		if err != nil {
			continue
		}

		outBytes, err := json.Marshal(socketmodels.OutChangeMessage{
			Type:   "CHANGE",
			Method: "UPDATE",
			Entity: "ROOM",
			Data:   string(jsonBytes),
		})
		if err != nil {
			continue
		}

		ss.SendDataToSubscription <- socketserver.SubscriptionDataMessage{
			Name: "room-display-data=" + changeEv.DocumentKey.ID.Hex(),
			Data: outBytes,
		}
	}
}

func watchRoomDeletes(db *mongo.Database, ss *socketserver.SocketServer) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic watching room deletes :", r)
		}
	}()
	cs, err := db.Collection("rooms").Watch(context.Background(), mongo.Pipeline{deletePipeline}, options.ChangeStream().SetFullDocument(options.UpdateLookup))
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
		id := changeEv.DocumentKey.ID

		db.Collection("room_internal_data").DeleteOne(context.Background(), bson.M{"_id": id})
		db.Collection("room_external_data").DeleteOne(context.Background(), bson.M{"_id": id})
		db.Collection("room_image").DeleteOne(context.Background(), bson.M{"_id": id})
		channelIds := []primitive.ObjectID{}
		if cursor, err := db.Collection("room_channels").Find(context.Background(), bson.M{"room_id": id}); err != nil {
			log.Fatal("CS CURSOR ERR : ", err)
			cursor.Close(context.Background())
		} else {
			for cursor.Next(context.Background()) {
				channel := &models.RoomChannel{}
				if err := cursor.Decode(&channel); err != nil {
					cursor.Close(context.Background())
					log.Fatal("CS CURSOR ERR : ", err)
				}
				channelIds = append(channelIds, channel.ID)
			}
			cursor.Close(context.Background())
		}
		db.Collection("room_channels").DeleteMany(context.Background(), bson.M{"room_id": id})
		db.Collection("room_channel_messages").DeleteMany(context.Background(), bson.M{"_id": bson.M{"$in": channelIds}})

		outBytes, err := json.Marshal(socketmodels.OutChangeMessage{
			Type:   "CHANGE",
			Method: "DELETE",
			Entity: "ROOM",
			Data:   `{"ID":"` + id.Hex() + `"}`,
		})
		if err != nil {
			continue
		}

		ss.SendDataToSubscription <- socketserver.SubscriptionDataMessage{
			Name: "room-display-data=" + changeEv.DocumentKey.ID.Hex(),
			Data: outBytes,
		}

		ss.DestroySubscription <- "room-display-data=" + id.Hex()
	}
}

func watchRoomChannelUpdates(db *mongo.Database, ss *socketserver.SocketServer) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic watching room channel updates :", r)
		}
	}()
	cs, err := db.Collection("room_channels").Watch(context.Background(), mongo.Pipeline{updatePipeline}, options.ChangeStream().SetFullDocument(options.UpdateLookup))
	if err != nil {
		log.Panicln("CS ERR : ", err.Error())
	}
	for cs.Next(context.Background()) {
		var changeEv struct {
			DocumentKey struct {
				ID primitive.ObjectID `bson:"_id"`
			} `bson:"documentKey"`
			FullDocument models.RoomChannel `bson:"fullDocument"`
		}
		err := cs.Decode(&changeEv)
		if err != nil {
			log.Println("CS DECODE ERROR : ", err)
			return
		}
		if changeEv.FullDocument.ToBeDeleted {
			db.Collection("room_internal_data").UpdateByID(context.Background(), changeEv.FullDocument.RoomID, bson.M{
				"$pull": bson.M{
					"channels": changeEv.DocumentKey.ID,
				},
			})
			db.Collection("room_channels").DeleteOne(context.Background(), bson.M{"_id": changeEv.FullDocument.ID})
		} else {
			outBytes, err := json.Marshal(changeEv.FullDocument)
			if err != nil {
				continue
			}

			ss.SendDataToSubscription <- socketserver.SubscriptionDataMessage{
				Name: "room-channel-data=" + changeEv.DocumentKey.ID.Hex(),
				Data: outBytes,
			}
		}
	}
}
