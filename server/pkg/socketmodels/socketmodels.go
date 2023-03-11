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

/* -------- ROOM MODELS -------- */

// TYPE: ROOM_OPEN_CHANNEL/ROOM_EXIT_CHANNEL
type RoomOpenExitChannel struct {
	Type    string `json:"TYPE"`
	Channel string `json:"channel"`
}

// TYPE: ROOM_MESSAGE
type RoomMessage struct {
	Type          string `json:"TYPE"`
	Content       string `json:"content"`
	Channel       string `json:"channel"`
	HasAttachment bool   `json:"has_attachment"`
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
	Type          string `json:"TYPE"`
	Content       string `json:"content"`
	ID            string `json:"ID"`
	Author        string `json:"author"`
	HasAttachment bool   `json:"has_attachment"`
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

/* -------- DIRECT MESSAGE, FRIEND REQUEST & INVITATION MODELS -------- */

// TYPE: DIRECT_MESSAGE
type DirectMessage struct {
	Type          string `json:"TYPE"`
	Content       string `json:"content"`
	Recipient     string `json:"recipient"`
	HasAttachment bool   `json:"has_attachment"`
}

// TYPE: ROOM_INVITATION
type InviteToRoom struct {
	Type      string `json:"TYPE"`
	Recipient string `json:"recipient"`
	RoomID    string `json:"room_id"`
}

// TYPE: FRIEND_REQUEST
type FriendRequest struct {
	Type      string `json:"TYPE"`
	Recipient string `json:"recipient"`
}

// TYPE: FRIEND_REQUEST_RESPONSE
type FriendRequestResponse struct {
	Type      string `json:"TYPE"`
	ID        string `json:"ID"`
	Accept    bool   `json:"accept"`
	Author    string `json:"author"`
	Recipient string `json:"recipient"`
}

// TYPE: ROOM_INVITATION_RESPONSE
type RoomInvitationResponse struct {
	Type   string `json:"TYPE"`
	ID     string `json:"ID"`
	Accept bool   `json:"accept"`
}

// TYPE: ROOM_INVITATION_DELETE
type RoomInvitationDelete struct {
	Type      string `json:"TYPE"`
	ID        string `json:"ID"`
	Recipient string `json:"recipient"`
}

// TYPE: DIRECT_MESSAGE_UPDATE
type DirectMessageUpdate struct {
	Type      string `json:"TYPE"`
	Content   string `json:"content"`
	Recipient string `json:"recipient"`
	ID        string `json:"ID"`
}

// TYPE: DIRECT_MESSAGE_DELETE
type DirectMessageDelete struct {
	Type      string `json:"TYPE"`
	Recipient string `json:"recipient"`
	ID        string `json:"ID"`
}

// TYPE: OUT_DIRECT_MESSAGE (no "TYPE" needed in model)
type OutDirectMessage struct {
	Content       string `json:"content"`
	ID            string `json:"ID"`
	Author        string `json:"author"`
	Recipient     string `json:"recipient"`
	HasAttachment bool   `json:"has_attachment"`
}

// TYPE: OUT_DIRECT_MESSAGE_UPDATE
type OutDirectMessageUpdate struct {
	Type      string `json:"TYPE"`
	Content   string `json:"content"`
	ID        string `json:"ID"`
	Author    string `json:"author"`
	Recipient string `json:"recipient"`
}

// TYPE: OUT_DIRECT_MESSAGE_DELETE
type OutDirectMessageDelete struct {
	Type      string `json:"TYPE"`
	ID        string `json:"ID"`
	Author    string `json:"author"`
	Recipient string `json:"recipient"`
}

// TYPE: OUT_INVITE
type OutInvite struct {
	Type      string `json:"TYPE"`
	ID        string `json:"ID"`
	Author    string `json:"author"`
	Recipient string `json:"recipient"`
	RoomID    string `json:"room_id"`
}

// TYPE: OUT_ROOM_INVITATION_DELETE
type OutRoomInvitationDelete struct {
	Type      string `json:"TYPE"`
	ID        string `json:"ID"`
	Author    string `json:"author"`
	Recipient string `json:"recipient"`
}

// TYPE: OUT_ROOM_INVITATION_RESPONSE
type OutRoomInvitationResponse struct {
	Type      string `json:"TYPE"`
	ID        string `json:"ID"`
	Author    string `json:"author"`
	Recipient string `json:"recipient"`
	Accept    bool   `json:"accept"`
}

// TYPE: OUT_FRIEND_REQUEST
type OutFriendRequest struct {
	Type      string `json:"TYPE"`
	ID        string `json:"ID"`
	Author    string `json:"author"`
	Recipient string `json:"recipient"`
}

// TYPE: OUT_FRIEND_REQUEST_DELETE
type OutFriendRequestDelete struct {
	Type      string `json:"TYPE"`
	ID        string `json:"ID"`
	Author    string `json:"author"`
	Recipient string `json:"recipient"`
}

// TYPE: OUT_FRIEND_REQUEST_RESPONSE
type OutFriendRequestResponse struct {
	Type      string `json:"TYPE"`
	ID        string `json:"ID"`
	Accept    bool   `json:"accept"`
	Author    string `json:"author"`
	Recipient string `json:"recipient"`
}

/* -------- BLOCK & UNBLOCK / BAN & UNBAN MODELS -------- */

// TYPE: BLOCK_USER/UNBLOCK_USER
type BlockUnblockUser struct {
	Type string `json:"TYPE"`
	UID  string `json:"uid"`
}

// TYPE: BAN_USER/UNBAN_USER
type BanUnbanUser struct {
	Type   string `json:"TYPE"`
	UID    string `json:"uid"`
	RoomID string `json:"room_id"`
}

/* -------- MISC -------- */

// TYPE: CHANGE
type OutChangeMessage struct {
	Type   string `json:"TYPE"`
	Method string `json:"METHOD"`
	Data   string `json:"DATA"`
	Entity string `json:"ENTITY"`
}

// TYPE: ATTACHMENT_PROGRESS (no "TYPE" needed in model)
type AttachmentProgress struct {
	MsgID  string  `json:"ID"`
	Ratio  float32 `json:"ratio"`
	Failed bool    `json:"err"`
}

// TYPE: ATTACHMENT_REQUEST (no "TYPE" needed in model)
type AttachmentRequest struct {
	MsgID  string `json:"ID"`
	IsRoom bool   `json:"is_room"`
}
