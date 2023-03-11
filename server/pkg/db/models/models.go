package models

import "go.mongodb.org/mongo-driver/bson/primitive"

/*
A lot of things are split up into seperate collections, mainly
so that changestreams don't get triggered by every minor change,
such as messages being added. Also because it's not necessary to send
all associated data on every request.

IDs for things like "UserMessagingData", "Pfp" and so on will be the
same as the document it refers to. For example the ID of a RoomImage
will be the same ID as the room it is for.
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

type DirectMessage struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"ID"`
	Content       string             `bson:"content,maxlength=300" json:"content"`
	CreatedAt     primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt     primitive.DateTime `bson:"updated_at" json:"updated_at"`
	Author        primitive.ObjectID `bson:"author" json:"author"`
	HasAttachment bool               `bson:"has_attachment" json:"has_attachment"`
}

type Invitation struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"ID"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
	Author    primitive.ObjectID `bson:"author" json:"author"`
	Recipient primitive.ObjectID `bson:"-" json:"recipient"`
	RoomID    primitive.ObjectID `bson:"room_id" json:"room_id"`
	Accepted  bool               `bson:"accepted" json:"accepted"`
	Declined  bool               `bson:"declined" json:"declined"`
}

type FriendRequest struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"ID"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
	Author    primitive.ObjectID `bson:"author" json:"author"`
	Recipient primitive.ObjectID `bson:"-" json:"recipient"`
	Accepted  bool               `bson:"accepted" json:"accepted"`
	Declined  bool               `bson:"declined" json:"declined"`
}

type UserMessagingData struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"ID"`
	Messages       []DirectMessage    `bson:"messages" json:"messages"`
	Invitations    []Invitation       `bson:"invitations" json:"invitations"`
	FriendRequests []FriendRequest    `bson:"friend_requests" json:"friend_requests"`
	// also includes invitations & friend requests
	MessagesSentTo []primitive.ObjectID `bson:"messages_sent_to" json:"messages_sent_to"`
	// also includes invitations & friend requests
	MessagesReceivedFrom []primitive.ObjectID `bson:"messages_received_from" json:"messages_received_from"`
	Blocked              []primitive.ObjectID `bson:"blocked" json:"blocked"`
	Friends              []primitive.ObjectID `bson:"friends" json:"friends"`
}

// Changes to pfp docs triggers changestream events
type Pfp struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"ID"`
	Binary primitive.Binary   `bson:"binary" json:"binary"`
}

/*---------------- Room structs ----------------*/

type RoomChannelMessage struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"ID"`
	Content       string             `bson:"content,maxlength=300" json:"content"`
	CreatedAt     primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt     primitive.DateTime `bson:"updated_at" json:"updated_at"`
	Author        primitive.ObjectID `bson:"author" json:"author"`
	HasAttachment bool               `bson:"has_attachment" json:"has_attachment"`
}

type RoomChannelMessages struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"ID"`
	Messages []RoomChannelMessage `bson:"messages" json:"messages"`
}

// Changes to room channel docs triggers changestream events
type RoomChannel struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"ID"`
	RoomID   primitive.ObjectID   `bson:"room_id" json:"-"`
	Name     string               `bson:"name" json:"name"`
	Messages []RoomChannelMessage `bson:"-" json:"messages"`
	// Need to use this because changeStream delete events dont return full document
	ToBeDeleted bool `bson:"to_be_deleted" json:"-"`
}

// Changes to room docs triggers changestream events
type Room struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"ID"`
	Name   string             `bson:"name" json:"name,maxlength=16"`
	Author primitive.ObjectID `bson:"author" json:"author"`
	// blur will be an empty string if the room has no image
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

// these things are seperate for a couple of reasons

type RoomInternalData struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"ID"`
	Channels    []primitive.ObjectID `bson:"channels" json:"channel"`
	MainChannel primitive.ObjectID   `bson:"main_channel" json:"main_channel"`
}

type RoomExternalData struct {
	ID      primitive.ObjectID   `bson:"_id,omitempty" json:"ID"`
	Private bool                 `bson:"private" json:"private"`
	Members []primitive.ObjectID `bson:"members" json:"members"`
	Banned  []primitive.ObjectID `bson:"banned" json:"banned"`
}

/*---------------- Attachment structs ----------------*/

type AttachmentChunk struct {
	// First chunk ID will be message ID
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"ID"`
	Data primitive.Binary   `bson:"data" json:"data"`
	// If its the last chunk this will be nil object ID
	NextChunkID primitive.ObjectID `bson:"next_chunk_id" json:"next_chunk_id"`
}
type AttachmentData struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"ID"`
	Meta   string             `bson:"meta" json:"meta"`
	Name   string             `bson:"name" json:"name"`
	Size   int                `bson:"size" json:"size"`
	Ratio  float32            `bson:"ratio" json:"ratio"`
	Failed bool               `bson:"failed" json:"failed"`
}
