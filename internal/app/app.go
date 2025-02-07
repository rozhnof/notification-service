package app

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	"github.com/rozhnof/notification-service/internal/application/services"
	"github.com/rozhnof/notification-service/internal/infrastructure/kafka"
	"github.com/rozhnof/notification-service/internal/pkg/mail"
)

func Run(ctx context.Context, mailSender mail.Sender, logger *slog.Logger, brokerList []string) error {
	notificationSender := services.NewEmailNotificationService(logger, mailSender)

	loginConsumerGroup, err := kafka.NewConsumerGroup(brokerList, "logins-group", []string{"logins"})
	if err != nil {
		return errors.Wrap(err, "failed init consumer group")
	}
	defer func() {
		if err := loginConsumerGroup.Close(); err != nil {
			logger.Error("failed close login consumer group")
		}
	}()
	go StartConsume(ctx, logger, loginConsumerGroup, notificationSender.SendLoginNotification)

	registerConsumerGroup, err := kafka.NewConsumerGroup(brokerList, "registers-group", []string{"registers"})
	if err != nil {
		return errors.Wrap(err, "failed init consumer group")
	}
	defer func() {
		if err := registerConsumerGroup.Close(); err != nil {
			logger.Error("failed close register consumer group")
		}
	}()
	go StartConsume(ctx, logger, registerConsumerGroup, notificationSender.SendRegisterMessage)

	<-ctx.Done()

	return nil
}

func StartConsume[T any](ctx context.Context, logger *slog.Logger, consumerGroup kafka.ConsumerGroup, msgHandler func(context.Context, T) error) {
	handler := kafka.NewConsumerGroupHandler(func(cm *sarama.ConsumerMessage) {
		var msg T

		if err := json.Unmarshal(cm.Value, &msg); err != nil {
			logger.Error("failed parse message", slog.String("error", err.Error()))
			return
		}

		if err := msgHandler(ctx, msg); err != nil {
			logger.Error("failed send message", slog.String("error", err.Error()))
			return
		}
	})

	go func() {
		for err := range consumerGroup.Errors() {
			logger.Warn("consumer group error", slog.String("error", err.Error()))
		}
	}()

	for {
		select {
		case <-ctx.Done():
			logger.Info("consumer group stopped")
			return
		default:
			if err := consumerGroup.Consume(ctx, &handler); err != nil {
				logger.Error("consume error", slog.String("error", err.Error()))
			}
		}
	}
}
