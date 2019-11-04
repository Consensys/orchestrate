package sarama

import (
	"github.com/Shopify/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

// NewClient creates a new sarama client and connects to one of the given broker addresses
func NewClient(addrs []string, conf *sarama.Config) (sarama.Client, error) {
	if err := conf.Validate(); err != nil {
		return nil, errors.ConfigError(err.Error()).SetComponent(component)
	}

	client, err := sarama.NewClient(addrs, conf)
	if err != nil {
		return nil, errors.KafkaConnectionError(err.Error()).SetComponent(component)
	}

	return client, nil
}
