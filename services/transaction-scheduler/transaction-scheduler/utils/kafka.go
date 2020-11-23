package utils

import (
	"context"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/handlers/multitenancy"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/transaction-scheduler/parsers"
)

func SendJobMessage(ctx context.Context, jobModel *models.Job, kafkaProducer sarama.SyncProducer, topic string) (partition int32, offset int64, err error) {
	log.WithContext(ctx).Debug("sending kafka message")

	txEnvelope := parsers.NewEnvelopeFromJobModel(jobModel, map[string]string{
		multitenancy.TenantIDMetadata: jobModel.Schedule.TenantID,
	})

	evlp, err := txEnvelope.Envelope()
	if err != nil {
		errMessage := "failed to craft envelope"
		log.WithContext(ctx).WithError(err).Error(errMessage)
		return 0, 0, errors.InvalidParameterError(errMessage)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
	}

	if partitionKey := evlp.PartitionKey(); partitionKey != "" {
		msg.Key = sarama.StringEncoder(partitionKey)
	}

	err = encoding.Marshal(txEnvelope, msg)
	if err != nil {
		errMessage := "failed to encode envelope"
		log.WithContext(ctx).WithError(err).Error(errMessage)
		return 0, 0, errors.InvalidParameterError(errMessage)
	}

	// Send message
	partition, offset, err = kafkaProducer.SendMessage(msg)
	if err != nil {
		errMessage := "could not produce kafka message"
		log.WithContext(ctx).WithError(err).Error(errMessage)
		return 0, 0, errors.KafkaConnectionError(errMessage)
	}

	log.WithField("envelope_id", txEnvelope.GetID()).
		WithField("job_type", evlp.GetJobTypeString()).
		Debug("envelope sent to kafka")
	return partition, offset, err
}
