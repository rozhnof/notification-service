package services

import (
	"github.com/IBM/sarama"
	"github.com/rozhnof/notification-service/internal/pkg/mail"
)

type MailSender struct {
	mailSender mail.Sender
}

func NewMailSender() *MailSender {
	return &MailSender{
		mailSender: mail.Sender{},
	}
}

func (s MailSender) SendLoginMessage(cm *sarama.ConsumerMessage) error {
	msg := mail.Message{
		Receiver: "",
		Subject:  "",
		Body:     "",
	}

	if err := s.mailSender.SendMessage(msg); err != nil {
		return err
	}

	return nil
}
