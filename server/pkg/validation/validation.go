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
type InsertRoomChannelData struct {
	Name          string `json:"name" validate:"required,max=24"`
	PromoteToMain bool   `json:"promote_to_main"`
}
type UpdateRoomChannelsData struct {
	UpdateData    []UpdateRoomChannelData `json:"update_data"`
	InsertData    []InsertRoomChannelData `json:"insert_data"`
	Delete        []string                `json:"delete_ids"`
	PromoteToMain string                  `json:"promote_to_main"`
}

type AttachmentMetadata struct {
	ID       string `json:"ID"`
	MimeType string `json:"mime_type"`
	Name     string `json:"name"`
	Size     int    `json:"size"`
}
