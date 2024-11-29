package mail

import (
	"fmt"
	"net/smtp"
)

type Sender struct {
	auth  smtp.Auth
	email string
	host  string
}

func NewSender(email string, password string, address string, port string) (*Sender, error) {
	var (
		auth = smtp.PlainAuth("", email, password, address)
		host = fmt.Sprintf("%s:%s", address, port)
	)

	return &Sender{
		auth:  auth,
		email: email,
		host:  host,
	}, nil
}

func (s *Sender) SendMessage(message Message) error {
	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		s.email,
		message.Receiver,
		message.Subject,
		message.Body,
	)

	if err := smtp.SendMail(s.host, s.auth, s.email, []string{message.Receiver}, []byte(msg)); err != nil {
		return err
	}

	return nil
}
