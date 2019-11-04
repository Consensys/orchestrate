package sarama

import (
	"context"

	"github.com/Shopify/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

type consumerGroup struct {
	g sarama.ConsumerGroup

	errors chan error
}

// NewConsumerGroupFromClient creates a new consumer group using the given client
func NewConsumerGroupFromClient(groupID string, client sarama.Client) (sarama.ConsumerGroup, error) {
	g, err := sarama.NewConsumerGroupFromClient(groupID, client)
	if err != nil {
		return nil, errors.KafkaConnectionError(err.Error()).SetComponent(component)
	}

	cg := &consumerGroup{
		g:      g,
		errors: make(chan error, client.Config().ChannelBufferSize),
	}

	// Pipe errors
	go func() {
		for err := range g.Errors() {
			cg.errors <- errors.KafkaConnectionError(err.Error()).SetComponent(component)
		}
	}()

	return cg, nil
}

// Consume implements ConsumerGroup.
func (c *consumerGroup) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	err := c.g.Consume(ctx, topics, handler)
	if err != nil {
		return errors.KafkaConnectionError(err.Error()).SetComponent(component)
	}
	return nil
}

// Errors implements ConsumerGroup.
func (c *consumerGroup) Errors() <-chan error {
	return c.errors
}

// Close implements ConsumerGroup.
func (c *consumerGroup) Close() error {
	err := c.g.Close()
	if err != nil {
		return errors.KafkaConnectionError(err.Error()).SetComponent(component)
	}
	return nil
}
