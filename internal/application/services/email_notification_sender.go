package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/rozhnof/notification-service/internal/pkg/mail"
)

type EmailNotificationService struct {
	logger *slog.Logger
	sender mail.Sender
}

func NewEmailNotificationService(logger *slog.Logger, sender mail.Sender) *EmailNotificationService {
	return &EmailNotificationService{
		logger: logger,
		sender: sender,
	}
}

func (s EmailNotificationService) SendLoginNotification(ctx context.Context, msg LoginMessage) error {
	mailMsg := mail.Message{
		Receiver: msg.Email,
		Subject:  "Login",
		Body:     "Your are logging in",
	}

	if err := s.sender.SendMessage(ctx, mailMsg); err != nil {
		return err
	}

	return nil
}

func (s EmailNotificationService) SendRegisterMessage(ctx context.Context, msg RegisterMessage) error {
	mailMsg := mail.Message{
		Receiver: msg.Email,
		Subject:  "Register",
		Body:     fmt.Sprintf("You are successfully registered! To confirm your email, please visit the following link: %s", msg.ConfirmLink),
	}

	if err := s.sender.SendMessage(ctx, mailMsg); err != nil {
		return err
	}

	return nil
}
