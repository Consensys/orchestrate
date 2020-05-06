package jobs

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
)

//go:generate mockgen -source=start_job.go -destination=mocks/start_job.go -package=mocks

const startJobComponent = "use-cases.start-job"

type StartJobUseCase interface {
	Execute(ctx context.Context, jobUUID string) error
}

// startJobUseCase is a use case to start a transaction job
type startJobUseCase struct {
	jobDataAgent   store.JobAgent
	logDataAgent   store.LogAgent
	kafkaProducer  sarama.SyncProducer
	txCrafterTopic string
}

// NewStartJobUseCase creates a new StartJobUseCase
func NewStartJobUseCase(
	jobDataAgent store.JobAgent,
	logDataAgent store.LogAgent,
	kafkaProducer sarama.SyncProducer,
	txCrafterTopic string,
) StartJobUseCase {
	return &startJobUseCase{
		jobDataAgent:   jobDataAgent,
		logDataAgent:   logDataAgent,
		kafkaProducer:  kafkaProducer,
		txCrafterTopic: txCrafterTopic,
	}
}

// Execute validates and creates a new transaction job
func (uc *startJobUseCase) Execute(ctx context.Context, jobUUID string) error {
	log.WithContext(ctx).WithField("job_uuid", jobUUID).Debug("starting job")

	job, err := uc.jobDataAgent.FindOneByUUID(ctx, jobUUID)
	if err != nil {
		return errors.FromError(err).ExtendComponent(startJobComponent)
	}

	var method tx.Method
	switch job.Type {
	case types.JobConstantinopleTransaction:
		method = tx.Method_ETH_SENDRAWTRANSACTION
	default:
		method = tx.Method_ETH_SENDRAWTRANSACTION
	}

	txEnvelope := &tx.TxEnvelope{
		Msg: &tx.TxEnvelope_TxRequest{TxRequest: &tx.TxRequest{
			Headers: nil, // TODO: Add the JWT token here? https://pegasys1.atlassian.net/browse/PO-544
			Chain:   job.Schedule.ChainID,
			Method:  method,
			Params: &tx.Params{
				From:           job.Transaction.Sender,
				To:             job.Transaction.Recipient,
				Gas:            job.Transaction.GasLimit,
				GasPrice:       job.Transaction.GasPrice,
				Value:          job.Transaction.Value,
				Nonce:          job.Transaction.Nonce,
				Data:           job.Transaction.Data,
				Raw:            job.Transaction.Raw,
				PrivateFor:     job.Transaction.PrivateFor,
				PrivateFrom:    job.Transaction.PrivateFrom,
				PrivacyGroupId: job.Transaction.PrivacyGroupID,
			},
			ContextLabels: job.Labels,
		}},
	}
	partition, offset, err := uc.sendMessage(ctx, txEnvelope)
	if err != nil {
		return errors.FromError(err).ExtendComponent(startJobComponent)
	}

	err = uc.logDataAgent.Insert(ctx, &models.Log{
		JobID:   job.ID,
		Status:  types.LogStatusStarted,
		Message: fmt.Sprintf("message sent to partition %v, offset %v and topic %v", partition, offset, uc.txCrafterTopic),
	})
	if err != nil {
		return errors.FromError(err).ExtendComponent(startJobComponent)
	}

	log.WithContext(ctx).WithField("job_uuid", jobUUID).Info("job started successfully")
	return nil
}

func (uc *startJobUseCase) sendMessage(ctx context.Context, txEnvelope *tx.TxEnvelope) (partition int32, offset int64, err error) {
	log.WithContext(ctx).Debug("sending kafka message")

	envelopeBytes, err := json.Marshal(txEnvelope)
	if err != nil {
		errMessage := "failed to encode envelope"
		log.WithContext(ctx).WithError(err).Error(errMessage)
		return 0, 0, errors.InvalidParameterError(errMessage)
	}

	msg := &sarama.ProducerMessage{
		Topic:   uc.txCrafterTopic,
		Value:   sarama.ByteEncoder(envelopeBytes),
		Headers: nil, // TODO: Add the JWT token here? https://pegasys1.atlassian.net/browse/PO-544
	}

	// Send message
	partition, offset, err = uc.kafkaProducer.SendMessage(msg)
	if err != nil {
		errMessage := "could not produce kafka message"
		log.WithContext(ctx).WithError(err).Error(errMessage)
		return 0, 0, errors.KafkaConnectionError(errMessage).ExtendComponent(startJobComponent)
	}

	return partition, offset, err
}
