package kafka

import (
	"context"
	"time"

	"github.com/IBM/sarama"
)

type ConsumerGroup struct {
	sarama.ConsumerGroup
	topics []string
}

func NewConsumerGroup(brokerList []string, groupID string, topics []string) (ConsumerGroup, error) {
	cfg := sarama.NewConfig()

	cfg.Version = sarama.MaxVersion
	cfg.Consumer.Return.Errors = true
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.Consumer.Group.Heartbeat.Interval = 3 * time.Second
	cfg.Consumer.Group.Session.Timeout = 60 * time.Second
	cfg.Consumer.Group.Rebalance.Timeout = 60 * time.Second
	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}

	consumerGroup, err := sarama.NewConsumerGroup(brokerList, groupID, cfg)
	if err != nil {
		return ConsumerGroup{}, err
	}

	return ConsumerGroup{
		ConsumerGroup: consumerGroup,
		topics:        topics,
	}, nil
}

func (cg ConsumerGroup) Consume(ctx context.Context, handler sarama.ConsumerGroupHandler) error {
	return cg.ConsumerGroup.Consume(ctx, cg.topics, handler)
}
