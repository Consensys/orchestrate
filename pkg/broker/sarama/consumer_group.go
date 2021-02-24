package sarama

import (
	"context"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/Shopify/sarama"
)

//go:generate mockgen -source=consumer_group.go -destination=mock/mock.go -package=mock
type ConsumerGroupSession interface {
	Claims() map[string][]int32
	MemberID() string
	GenerationID() int32
	MarkOffset(topic string, partition int32, offset int64, metadata string)
	Commit()
	ResetOffset(topic string, partition int32, offset int64, metadata string)
	MarkMessage(msg *sarama.ConsumerMessage, metadata string)
	Context() context.Context
}

type ConsumerGroupClaim interface {
	Topic() string
	Partition() int32
	InitialOffset() int64
	HighWaterMarkOffset() int64
	Messages() <-chan *sarama.ConsumerMessage
}

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
