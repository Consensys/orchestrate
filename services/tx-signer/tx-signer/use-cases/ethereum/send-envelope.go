package ethereum

import (
	"context"

	"github.com/golang/protobuf/proto"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-signer/tx-signer/use-cases"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

const sendEnvelopeComponent = "use-cases.send-envelope"

// sendEnvelopeUseCase is a use case to sign a public Ethereum transaction
type sendEnvelopeUseCase struct {
	producer sarama.SyncProducer
}

// NewSendEnvelopeUseCase creates a new SendEnvelopeUseCase
func NewSendEnvelopeUseCase(producer sarama.SyncProducer) usecases.SendEnvelopeUseCase {
	return &sendEnvelopeUseCase{
		producer: producer,
	}
}

// Execute sends an envelope to the tx-sender
func (uc *sendEnvelopeUseCase) Execute(ctx context.Context, protoMessage proto.Message, topic, partitionKey string) error {
	logger := log.WithContext(ctx).WithField("topic", topic)
	logger.Debug("sending envelope")

	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	// Set key for Kafka partitions
	if partitionKey != "" {
		msg.Key = sarama.StringEncoder(partitionKey)
	}

	b, err := encoding.Marshal(protoMessage)
	if err != nil {
		errMessage := "failed to marshal envelope as request"
		log.WithError(err).Error(errMessage)
		return errors.EncodingError(errMessage).ExtendComponent(sendEnvelopeComponent)
	}
	msg.Value = sarama.ByteEncoder(b)

	partition, offset, err := uc.producer.SendMessage(msg)
	if err != nil {
		errMessage := "failed to produce kafka message"
		log.WithError(err).Error(errMessage)
		return errors.KafkaConnectionError(errMessage).ExtendComponent(sendEnvelopeComponent)
	}

	log.WithField("partition", partition).
		WithField("offset", offset).
		WithField("topic", msg.Topic).
		Debug("envelope successfully sent")
	return nil
}
