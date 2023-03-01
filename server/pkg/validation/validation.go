package validation

type Credentials struct {
	Username string `json:"username" validate:"required,min=2,max=16"`
	Password string `json:"password" validate:"required,min=2,max=100"`
}

type Room struct {
	Name    string `json:"name" validate:"required,min=2,max=16"`
	Private bool   `json:"private"`
}

type UserSearch struct {
	Username string `json:"username" validate:"max=16"`
}
