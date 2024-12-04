package app

import (
	"fmt"

	"github.com/rozhnof/notification-service/internal/domain/models"
	"github.com/rozhnof/notification-service/internal/pkg/mail"
)

type MailSender struct {
	sender *mail.Sender
}

func NewMailSender(sender *mail.Sender) MailSender {
	return MailSender{
		sender: sender,
	}
}

func (s MailSender) SendLoginMessage(msg models.LoginMessage) error {
	mailMsg := mail.Message{
		Receiver: msg.Email,
		Subject:  "Login",
		Body:     "Your are logging in",
	}

	if err := s.sender.SendMessage(mailMsg); err != nil {
		return err
	}

	return nil
}

func (s MailSender) SendRegisterMessage(msg models.RegisterMessage) error {
	mailMsg := mail.Message{
		Receiver: msg.Email,
		Subject:  "Register",
		Body:     fmt.Sprintf("You are successfully registered! To confirm your email, please visit the following link: %s", msg.ConfirmLink),
	}

	if err := s.sender.SendMessage(mailMsg); err != nil {
		return err
	}

	return nil
}
