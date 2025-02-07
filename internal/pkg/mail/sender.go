package mail

import (
	"context"
	"fmt"
	"net/smtp"
)

type Sender struct {
	auth  smtp.Auth
	email string
	host  string
}

func NewSender(email string, password string, address string, port string) (Sender, error) {
	var (
		auth = smtp.PlainAuth("", email, password, address)
		host = fmt.Sprintf("%s:%s", address, port)
	)

	return Sender{
		auth:  auth,
		email: email,
		host:  host,
	}, nil
}

func (s *Sender) SendMessage(_ context.Context, msg Message) error {
	msgBytes := s.buildMessage(msg)

	if err := smtp.SendMail(s.host, s.auth, s.email, []string{msg.Receiver}, msgBytes); err != nil {
		return err
	}

	return nil
}

func (s *Sender) buildMessage(msg Message) []byte {
	return []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		s.email,
		msg.Receiver,
		msg.Subject,
		msg.Body,
	))
}
