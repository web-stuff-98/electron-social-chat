package models

import "go.mongodb.org/mongo-driver/bson/primitive"

/*---------------- User structs (session in redis) ----------------*/

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"ID"`
	Username  string             `bson:"username,maxlength=16" json:"username"`
	Password  string             `bson:"password" json:"-"`
	Base64pfp string             `bson:"-" json:"base64pfp,omitempty"`
	IsOnline  bool               `bson:"-" json:"online"`
}

type Pfp struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"ID"`
	Binary primitive.Binary   `bson:"binary" json:"binary"`
}

/*---------------- Room structs ----------------*/

type RoomChannelMessage struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"ID"`
	Content   string             `bson:"content,maxlength=200" json:"content"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at" json:"updated_at"`
	Author    primitive.ObjectID `bson:"author" json:"author"`
}

type RoomChannelMessages struct {
	Messages []RoomChannelMessage `bson:"messages" json:"messages"`
}

type RoomChannel struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"ID"`
	Name string             `bson:"name" json:"name"`
}

type Room struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"ID"`
	Name   string             `bson:"name" json:"name,maxlength=16"`
	Author primitive.ObjectID `bson:"author" json:"author"`
}

type RoomImage struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Binary primitive.Binary   `bson:"binary"`
}

// Potentially heavier data (room might have a lot of messages) that should be loaded after the room is entered
type RoomInternalData struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"ID"`
	// Channels also contain all the messages, so potentially heavy.
	Channels    []RoomChannel      `bson:"channels" json:"channel"`
	MainChannel primitive.ObjectID `bson:"main_channel" json:"main_channel"`
}

// Lighter data that will be loaded before the room is opened
type RoomExternalData struct {
	ID      primitive.ObjectID   `bson:"_id,omitempty" json:"ID"`
	Private bool                 `bson:"private" json:"private"`
	Members []primitive.ObjectID `bson:"members" json:"members"`
	Banned  []primitive.ObjectID `bson:"banned" json:"banned"`
}
