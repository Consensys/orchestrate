package sarama

import (
	"github.com/Shopify/sarama"
	errors "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
)

// NewSyncProducerFromClient creates a new sarama.SyncProducer using the given client
func NewSyncProducerFromClient(client sarama.Client) (sarama.SyncProducer, error) {
	p, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, errors.KafkaConnectionError(err.Error()).SetComponent(component)
	}
	return p, nil
}
