package kafka

import (
	"github.com/IBM/sarama"
)

type ConsumerGroupHandler struct {
	ready   chan struct{}
	handler func(*sarama.ConsumerMessage)
}

func NewConsumerGroupHandler(handler func(*sarama.ConsumerMessage)) ConsumerGroupHandler {
	return ConsumerGroupHandler{
		ready:   make(chan struct{}),
		handler: handler,
	}
}

func (c *ConsumerGroupHandler) Ready() <-chan struct{} {
	return c.ready
}

func (c *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	close(c.ready)

	return nil
}

func (c *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		c.handler(message)
		session.MarkMessage(message, "")
	}

	return nil
}
