package services

import (
	"encoding/json"
	"log/slog"

	"github.com/rozhnof/notification-service/internal/domain/models"
)

type Sender interface {
	SendLoginMessage(msg models.LoginMessage) error
	SendRegisterMessage(msg models.RegisterMessage) error
}

type NotificationSender struct {
	logger  *slog.Logger
	senders []Sender
}

func NewNotificationSender(senders []Sender, logger *slog.Logger) *NotificationSender {
	return &NotificationSender{
		senders: senders,
		logger:  logger,
	}
}

func (s NotificationSender) ConsumeLoginMessage(message []byte) error {
	var loginMessage models.LoginMessage

	if err := json.Unmarshal(message, &loginMessage); err != nil {
		return err
	}

	for _, sender := range s.senders {
		if err := sender.SendLoginMessage(loginMessage); err != nil {
			return err
		}
	}

	return nil
}

func (s NotificationSender) ConsumeRegisterMessage(message []byte) error {
	var registerMessage models.RegisterMessage

	if err := json.Unmarshal(message, &registerMessage); err != nil {
		return err
	}

	for _, sender := range s.senders {
		if err := sender.SendRegisterMessage(registerMessage); err != nil {
			return err
		}
	}

	return nil
}
