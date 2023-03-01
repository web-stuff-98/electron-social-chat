package db

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Collections struct {
	UserCollection *mongo.Collection
	PfpCollection  *mongo.Collection

	RoomCollection             *mongo.Collection
	RoomInternalDataCollection *mongo.Collection
	RoomExternalDataCollection *mongo.Collection
	RoomImageCollection        *mongo.Collection
}

func Init() (*mongo.Database, *Collections) {
	log.Println("Connecting to MongoDB")

	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatal(err)
	}
	client.Connect(context.Background())
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	DB := client.Database(os.Getenv("MONGODB_DB"))

	colls := &Collections{
		UserCollection: DB.Collection("users"),
		PfpCollection:  DB.Collection("pfps"),

		RoomCollection:             DB.Collection("rooms"),
		RoomInternalDataCollection: DB.Collection("room_internal_data"),
		RoomExternalDataCollection: DB.Collection("room_external_data"),
		RoomImageCollection:        DB.Collection("room_image"),
	}

	colls.UserCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.M{
			"username": "text",
		},
		Options: options.Index().SetName("username_text"),
	})

	log.Println("Connected to MongoDB")

	return DB, colls
}
