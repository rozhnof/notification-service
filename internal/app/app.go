package app

import (
	"context"
	"log/slog"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	"github.com/rozhnof/notification-service/internal/application/services"
	"github.com/rozhnof/notification-service/internal/infrastructure/kafka"
	"github.com/rozhnof/notification-service/internal/pkg/mail"
)

func Run(ctx context.Context, mailSender *mail.Sender, logger *slog.Logger, brokerList []string) error {
	senders := []services.Sender{
		NewMailSender(mailSender),
	}

	notificationSender := services.NewNotificationSender(senders, logger)

	loginConsumerGroup, err := kafka.NewConsumerGroup(brokerList, "logins-group")
	if err != nil {
		return errors.Wrap(err, "failed init consumer group")
	}
	defer func() {
		if err := loginConsumerGroup.Close(); err != nil {
			logger.Error("failed close login consumer group")
		}
	}()
	go StartConsume(ctx, logger, []string{"logins"}, loginConsumerGroup, notificationSender.ConsumeLoginMessage)

	registerConsumerGroup, err := kafka.NewConsumerGroup(brokerList, "registers-group")
	if err != nil {
		return errors.Wrap(err, "failed init consumer group")
	}
	defer func() {
		if err := registerConsumerGroup.Close(); err != nil {
			logger.Error("failed close register consumer group")
		}
	}()
	go StartConsume(ctx, logger, []string{"registers"}, registerConsumerGroup, notificationSender.ConsumeRegisterMessage)

	<-ctx.Done()

	return nil
}

func StartConsume(ctx context.Context, logger *slog.Logger, topics []string, consumerGroup kafka.ConsumerGroup, messageHandler func([]byte) error) {
	handler := kafka.NewConsumerGroupHandler(func(cm *sarama.ConsumerMessage) {
		if err := messageHandler(cm.Value); err != nil {
			logger.Error("failed send email message", slog.String("error", err.Error()))
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
			if err := consumerGroup.Consume(ctx, topics, &handler); err != nil {
				logger.Error("consumer group: consume error", slog.String("error", err.Error()))
			}
		}
	}
}
