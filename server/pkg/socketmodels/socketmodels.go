package socketmodels

/*
	Models for messages sent through the websocket, encoded into []bytes from json marshal

	When a socket message is sent out the "event type" is keyed as TYPE, when a socket message
	is recieved on the server it should be keyed as event_type, this is just so that its a bit
	easier to tell which models for sending data out, and which are for receiving data from
	the client.
*/

// TYPE: OPEN_SUBSCRIPTION/CLOSE_SUBSCRIPTION
type OpenCloseSubscription struct {
	Name string `json:"name"`
}

// TYPE: OPEN_SUBSCRIPTIONS
type OpenCloseSubscriptions struct {
	Names []string `json:"names"`
}

// TYPE: ROOM_MESSAGE/ROOM_MESSAGE_DELETE/ROOM_MESSAGE_UPDATE/PRIVATE_MESSAGE/PRIVATE_MESSAGE_DELETE/PRIVATE_MESSAGE_UPDATE/PRIVATE_MESSAGE_INVITE_RESPONDED/POST_VOTE/POST_COMMENT_VOTE/ATTACHMENT_PROGRESS/ATTACHMENT_COMPLETE/RESPONSE_MESSAGE/NOTIFICATIONS
type OutMessage struct {
	Type string `json:"TYPE"`
	Data string `json:"DATA"`
}

// TYPE: CHANGE
type OutChangeMessage struct {
	Type   string `json:"TYPE"`
	Method string `json:"METHOD"`
	Data   string `json:"DATA"`
	Entity string `json:"ENTITY"`
}
