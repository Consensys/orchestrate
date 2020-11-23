package sarama

import (
	"github.com/Shopify/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
)

type syncProducer struct {
	p sarama.SyncProducer
}

// NewSyncProducerFromClient creates a new sarama.SyncProducer using the given client
func NewSyncProducerFromClient(client sarama.Client) (sarama.SyncProducer, error) {
	p, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, errors.KafkaConnectionError(err.Error()).SetComponent(component)
	}
	return &syncProducer{p: p}, nil
}

// SendMessage produces a given message
func (sp *syncProducer) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	partition, offset, err = sp.p.SendMessage(msg)
	if err != nil {
		return partition, offset, errors.KafkaConnectionError(err.Error()).SetComponent(component)
	}
	return partition, offset, nil
}

// SendMessages produces a given set of messages
func (sp *syncProducer) SendMessages(msgs []*sarama.ProducerMessage) error {
	err := sp.p.SendMessages(msgs)
	if err != nil {
		return errors.KafkaConnectionError(err.Error()).SetComponent(component)
	}
	return nil
}

// Close shuts down the producer and waits for any buffered messages to be flushed
func (sp *syncProducer) Close() error {
	err := sp.p.Close()
	if err != nil {
		return errors.KafkaConnectionError(err.Error()).SetComponent(component)
	}
	return nil
}
