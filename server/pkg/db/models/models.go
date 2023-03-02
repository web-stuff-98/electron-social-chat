package models

import "go.mongodb.org/mongo-driver/bson/primitive"

/*
A lot of things are split up into seperate collections, mainly
so that changestreams don't get triggered by every minor change,
such as messages being added. Changestreams should only be triggered
in certain cases.
*/

/*---------------- User structs (session in redis) ----------------*/

// Users cannot change their names (there is no reason for this, could add this)
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"ID"`
	Username  string             `bson:"username,maxlength=16" json:"username"`
	Password  string             `bson:"password" json:"-"`
	Base64pfp string             `bson:"-" json:"base64pfp,omitempty"`
	IsOnline  bool               `bson:"-" json:"online"`
}

// Changes to pfp docs triggers changestream events
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
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"ID"`
	Messages []RoomChannelMessage `bson:"messages" json:"messages"`
}

// Changes to room channel docs triggers changestream events
type RoomChannel struct {
	// ID for the corresponding RoomChannelMessages doc will be the same
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"ID"`
	RoomID   primitive.ObjectID   `bson:"room_id" json:"-"`
	Name     string               `bson:"name" json:"name"`
	Messages []RoomChannelMessage `bson:"-" json:"messages"`
}

// Changes to room docs triggers changestream events
type Room struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"ID"`
	Name   string             `bson:"name" json:"name,maxlength=16"`
	Author primitive.ObjectID `bson:"author" json:"author"`
	// blur will not be present if the room has no image
	Blur string `bson:"blur" json:"blur"`

	Private bool                 `bson:"-" json:"is_private"`
	Members []primitive.ObjectID `bson:"-" json:"members"`
	Banned  []primitive.ObjectID `bson:"-" json:"banned"`

	Channels    []primitive.ObjectID `bson:"-" json:"channels"`
	MainChannel primitive.ObjectID   `bson:"-" json:"main_channel"`
}

type RoomImage struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Binary primitive.Binary   `bson:"binary"`
}

// Will only be accessed when a client opens the room
type RoomInternalData struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"ID"`
	Channels    []primitive.ObjectID `bson:"channels" json:"channel"`
	MainChannel primitive.ObjectID   `bson:"main_channel" json:"main_channel"`
}

// Data that will be accessed before the room is opened
type RoomExternalData struct {
	ID      primitive.ObjectID   `bson:"_id,omitempty" json:"ID"`
	Private bool                 `bson:"private" json:"private"`
	Members []primitive.ObjectID `bson:"members" json:"members"`
	Banned  []primitive.ObjectID `bson:"banned" json:"banned"`
}
