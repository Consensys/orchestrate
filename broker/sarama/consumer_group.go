package sarama

import (
	"github.com/Shopify/sarama"
	errors "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
)

// NewConsumerGroupFromClient creates a new consumer group using the given client
func NewConsumerGroupFromClient(groupID string, client sarama.Client) (sarama.ConsumerGroup, error) {
	g, err := sarama.NewConsumerGroupFromClient(groupID, client)
	if err != nil {
		return nil, errors.KafkaConnectionError(err.Error()).SetComponent(component)
	}
	return g, nil
}
