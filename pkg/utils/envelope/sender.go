package envelope

import (
	"github.com/Shopify/sarama"
	encoding "github.com/consensys/orchestrate/pkg/encoding/sarama"
	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/utils"
)

func SendJobMessage(job *entities.Job, kafkaProducer sarama.SyncProducer, topic string) (partition int32, offset int64, err error) {
	txEnvelope := NewEnvelopeFromJob(job, map[string]string{
		utils.TenantIDMetadata: job.TenantID,
	})

	evlp, err := txEnvelope.Envelope()
	if err != nil {
		return 0, 0, errors.InvalidParameterError("failed to craft envelope")
	}

	logger := log.NewLogger()
	logger.WithField("txType", evlp.TransactionType).Warn("SendJobMessage")

	msg := &sarama.ProducerMessage{
		Topic: topic,
	}

	if partitionKey := evlp.PartitionKey(); partitionKey != "" {
		msg.Key = sarama.StringEncoder(partitionKey)
	}

	err = encoding.Marshal(txEnvelope, msg)
	if err != nil {
		return 0, 0, errors.InvalidParameterError("failed to encode envelope")
	}

	// Send message
	partition, offset, err = kafkaProducer.SendMessage(msg)
	if err != nil {
		return 0, 0, errors.KafkaConnectionError("could not produce kafka message")
	}

	return partition, offset, err
}
