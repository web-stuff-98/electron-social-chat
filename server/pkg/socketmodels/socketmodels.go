package socketmodels

/*
	Models for messages sent through the websocket, encoded into []bytes from json marshal

	When a socket message is sent out the "event type" is keyed as TYPE, when a socket message
	is recieved on the server it should be keyed as event_type, this is just so that its a bit
	easier to tell which models for sending data out, and which are for receiving data from
	the client.
*/

// TYPE: WATCH_USER/STOP_WATCHING_USER/WATCH_ROOM/STOP_WATCHING_ROOM
type WatchStopWatching struct {
	Type string `json:"TYPE"`
	ID   string `json:"ID"`
}

// TYPE: ROOM_OPEN_CHANNEL/ROOM_EXIT_CHANNEL
type RoomOpenExitChannel struct {
	Type    string `json:"TYPE"`
	Channel string `json:"channel"`
}

// TYPE: ROOM_MESSAGE
type RoomMessage struct {
	Type    string `json:"TYPE"`
	Content string `json:"content"`
	Channel string `json:"channel"`
}

// TYPE: ROOM_MESSAGE_UPDATE
type RoomMessageUpdate struct {
	Type    string `json:"TYPE"`
	Content string `json:"content"`
	Channel string `json:"channel"`
	ID      string `json:"ID"`
}

// TYPE: ROOM_MESSAGE_DELETE
type RoomMessageDelete struct {
	Type    string `json:"TYPE"`
	Channel string `json:"channel"`
	ID      string `json:"ID"`
}

// TYPE: OUT_ROOM_MESSAGE
type OutRoomMessage struct {
	Type    string `json:"TYPE"`
	Content string `json:"content"`
	ID      string `json:"ID"`
	Author  string `json:"author"`
}

// TYPE: OUT_ROOM_MESSAGE_UPDATE
type OutRoomMessageUpdate struct {
	Type    string `json:"TYPE"`
	Content string `json:"content"`
	ID      string `json:"ID"`
}

// TYPE: OUT_ROOM_MESSAGE_DELETE
type OutRoomMessageDelete struct {
	Type string `json:"TYPE"`
	ID   string `json:"ID"`
}

// TYPE: CHANGE
type OutChangeMessage struct {
	Type   string `json:"TYPE"`
	Method string `json:"METHOD"`
	Data   string `json:"DATA"`
	Entity string `json:"ENTITY"`
}
