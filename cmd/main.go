package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/rozhnof/notification-service/internal/app"
	"github.com/rozhnof/notification-service/internal/infrastructure/kafka"
	"github.com/rozhnof/notification-service/internal/pkg/config"
	"github.com/rozhnof/notification-service/internal/pkg/mail"
)

const (
	EnvConfigPath = "CONFIG_PATH"
)

type Config struct {
	Mode   string        `yaml:"mode"    env-required:"true"`
	Kafka  config.Kafka  `yaml:"kafka"   env-required:"true"`
	Logger config.Logger `yaml:"logging" env-required:"true"`
	Mail   config.Mail
}

type LoginMessage struct {
	Email string `json:"email"`
}

type RegisterMessage struct {
	Email       string `json:"email"`
	ConfirmLink string `json:"confirm_link"`
}

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	cfg, err := config.NewConfig[Config](os.Getenv(EnvConfigPath))
	if err != nil {
		log.Fatal(err)
	}

	logger, err := app.NewLogger(cfg.Logger)
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("init logger success")

	emailSender, err := mail.NewSender(cfg.Mail.Email, cfg.Mail.Password, cfg.Mail.Address, cfg.Mail.Port)
	if err != nil {
		logger.Error("failed init email sender")
		return
	}

	loginConsumerGroup, err := kafka.NewConsumerGroup(cfg.Kafka.BrokerList, "logins-group")
	if err != nil {
		logger.Error("failed init login consumer group")
		return
	}
	defer func() {
		if err := loginConsumerGroup.Close(); err != nil {
			logger.Error("failed close login consumer group")
		}
	}()
	logger.Info("init login consumer group success")

	registerConsumerGroup, err := kafka.NewConsumerGroup(cfg.Kafka.BrokerList, "registers-group")
	if err != nil {
		logger.Error("failed init register consumer group")
		return
	}
	defer func() {
		if err := registerConsumerGroup.Close(); err != nil {
			logger.Error("failed close register consumer group")
		}
	}()
	logger.Info("init register consumer group success")

	go StartLoginConsume(ctx, logger, loginConsumerGroup, emailSender)
	go StartRegisterConsume(ctx, logger, registerConsumerGroup, emailSender)

	<-ctx.Done()
}

func StartLoginConsume(ctx context.Context, logger *slog.Logger, consumerGroup kafka.ConsumerGroup, emailSender *mail.Sender) {
	go func() {
		for err := range consumerGroup.Errors() {
			logger.Warn("login consumer group error", slog.String("error", err.Error()))
		}
	}()

	consumerGroupHandler := kafka.NewConsumerGroupHandler(func(cm *sarama.ConsumerMessage) {
		var loginMessage LoginMessage

		if err := json.Unmarshal(cm.Value, &loginMessage); err != nil {
			logger.Error("failed unmarshal kafka message", slog.String("error", err.Error()))
			return
		}

		emailMessage := mail.Message{
			Receiver: loginMessage.Email,
			Subject:  "New login",
			Body:     "You are logging in",
		}

		if err := emailSender.SendMessage(emailMessage); err != nil {
			logger.Error("failed send email message", slog.String("error", err.Error()))
			return
		}
	})

	for {
		select {
		case <-ctx.Done():
			logger.Info("login consumer group stopped")
			return
		default:
			if err := consumerGroup.Consume(ctx, []string{"logins"}, &consumerGroupHandler); err != nil {
				logger.Error("login consumer group: consume error", slog.String("error", err.Error()))
			}
		}
	}
}

func StartRegisterConsume(ctx context.Context, logger *slog.Logger, consumerGroup kafka.ConsumerGroup, emailSender *mail.Sender) {
	go func() {
		for err := range consumerGroup.Errors() {
			logger.Warn("register consumer group error", slog.String("error", err.Error()))
		}
	}()

	consumerGroupHandler := kafka.NewConsumerGroupHandler(func(cm *sarama.ConsumerMessage) {
		var registerMessage RegisterMessage

		if err := json.Unmarshal(cm.Value, &registerMessage); err != nil {
			logger.Error("failed unmarshal kafka message", slog.String("error", err.Error()))
			return
		}

		emailMessage := mail.Message{
			Receiver: registerMessage.Email,
			Subject:  "Registration",
			Body:     fmt.Sprintf("You are successfully registered! To confirm your email, please visit the following link: %s", registerMessage.ConfirmLink),
		}

		if err := emailSender.SendMessage(emailMessage); err != nil {
			logger.Error("failed send email message", slog.String("error", err.Error()))
			return
		}
	})

	for {
		select {
		case <-ctx.Done():
			logger.Info("register consumer group stopped")
			return
		default:
			if err := consumerGroup.Consume(ctx, []string{"registers"}, &consumerGroupHandler); err != nil {
				logger.Error("register consumer group: consume error", slog.String("error", err.Error()))
			}
		}
	}
}
