package models

type RegisterMessage struct {
	Email       string `json:"email"`
	ConfirmLink string `json:"confirm_link"`
}
