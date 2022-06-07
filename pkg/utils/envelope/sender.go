package envelope

import (
	"encoding/json"

	"github.com/Shopify/sarama"
	encoding "github.com/consensys/orchestrate/pkg/encoding/sarama"
	"github.com/consensys/orchestrate/pkg/errors"
	authutils "github.com/consensys/orchestrate/pkg/toolkit/app/auth/utils"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/entities"
)

func SendJobMessage(job *entities.Job, kafkaProducer sarama.SyncProducer, topic string, userInfo *multitenancy.UserInfo) (partition int32, offset int64, err error) {
	bUserInfo, _ := json.Marshal(userInfo)
	txEnvelope := NewEnvelopeFromJob(job, map[string]string{
		authutils.UserInfoHeader: string(bUserInfo),
	})

	evlp, err := txEnvelope.Envelope()
	if err != nil {
		return 0, 0, errors.InvalidParameterError("failed to craft envelope (%s)", err.Error())
	}

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
