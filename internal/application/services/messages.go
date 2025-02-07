package services

type LoginMessage struct {
	Email string `json:"email"`
}

type RegisterMessage struct {
	Email       string `json:"email"`
	ConfirmLink string `json:"confirm_link"`
}
