package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/rozhnof/notification-service/internal/app"
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

	mailSender, err := mail.NewSender(cfg.Mail.Email, cfg.Mail.Password, cfg.Mail.Address, cfg.Mail.Port)
	if err != nil {
		logger.Error("failed init email sender")
		return
	}

	if err := app.Run(ctx, mailSender, logger, cfg.Kafka.BrokerList); err != nil {
		logger.Error("notification app error", slog.String("error", err.Error()))
		return
	}
	logger.Info("notification app stopped")
}
