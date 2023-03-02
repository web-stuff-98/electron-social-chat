package validation

type Credentials struct {
	Username string `json:"username" validate:"required,min=2,max=16"`
	Password string `json:"password" validate:"required,min=2,max=100"`
}

type Room struct {
	Name    string `json:"name" validate:"required,min=2,max=16"`
	Private bool   `json:"is_private"`
}

type UserSearch struct {
	Username string `json:"username" validate:"max=16"`
}

type UpdateRoomChannelData struct {
	ID   string `json:"ID"`
	Name string `json:"name"`
}

type UpdateRoomChannelsData struct {
	UpdateData    []UpdateRoomChannelData `bson:"update_data"`
	InsertData    []string                `bson:"insert_data"`
	Delete        []string                `bson:"delete_ids"`
	PromoteToMain string                  `bson:"promote_to_main"`
}
