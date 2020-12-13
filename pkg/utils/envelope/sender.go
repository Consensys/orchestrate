package envelope

import (
	"context"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
)

func SendJobMessage(ctx context.Context, job *entities.Job, kafkaProducer sarama.SyncProducer, topic string) (partition int32, offset int64, err error) {
	log.WithContext(ctx).Debug("sending kafka message")

	txEnvelope := NewEnvelopeFromJob(job, map[string]string{
		utils.TenantIDMetadata: job.TenantID,
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
		WithField("topic", topic).
		Debug("envelope successfully sent")
	return partition, offset, err
}
