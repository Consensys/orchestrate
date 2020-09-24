package parsers

import (
	"math/big"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
)

func NewJobModelFromEntities(job *entities.Job, scheduleID *int) *models.Job {
	jobModel := &models.Job{
		UUID:         job.UUID,
		ChainUUID:    job.ChainUUID,
		Type:         job.Type,
		NextJobUUID:  job.NextJobUUID,
		Labels:       job.Labels,
		InternalData: job.InternalData,
		ScheduleID:   scheduleID,
		Schedule: &models.Schedule{
			UUID:     job.ScheduleUUID,
			TenantID: job.TenantID,
		},
		Logs:      []*models.Log{},
		CreatedAt: job.CreatedAt,
		UpdatedAt: job.UpdatedAt,
	}

	if scheduleID != nil {
		jobModel.Schedule.ID = *scheduleID
	}

	if job.Transaction != nil {
		jobModel.Transaction = NewTransactionModelFromEntities(job.Transaction)
	}

	for _, log := range job.Logs {
		jobModel.Logs = append(jobModel.Logs, NewLogModelFromEntity(log))
	}

	return jobModel
}

func NewJobEntityFromModels(jobModel *models.Job) *entities.Job {
	job := &entities.Job{
		UUID:         jobModel.UUID,
		ChainUUID:    jobModel.ChainUUID,
		NextJobUUID:  jobModel.NextJobUUID,
		Type:         jobModel.Type,
		Labels:       jobModel.Labels,
		InternalData: jobModel.InternalData,
		Logs:         []*entities.Log{},
		CreatedAt:    jobModel.CreatedAt,
		UpdatedAt:    jobModel.UpdatedAt,
	}

	if jobModel.Schedule != nil {
		job.ScheduleUUID = jobModel.Schedule.UUID
		job.TenantID = jobModel.Schedule.TenantID
	}

	if jobModel.Transaction != nil {
		job.Transaction = NewTransactionEntityFromModels(jobModel.Transaction)
	}

	for _, logModel := range jobModel.Logs {
		job.Logs = append(job.Logs, NewLogEntityFromModels(logModel))
	}

	return job
}

func NewEnvelopeFromJobModel(job *models.Job, headers map[string]string) *tx.TxEnvelope {
	contextLabels := job.Labels
	if contextLabels == nil {
		contextLabels = map[string]string{}
	}
	contextLabels[tx.ScheduleUUIDLabel] = job.Schedule.UUID
	contextLabels[tx.NextJobUUIDLabel] = job.NextJobUUID
	contextLabels[tx.PriorityLabel] = job.InternalData.Priority
	contextLabels[tx.ParentJobUUIDLabel] = job.InternalData.ParentJobUUID

	txEnvelope := &tx.TxEnvelope{
		Msg: &tx.TxEnvelope_TxRequest{TxRequest: &tx.TxRequest{
			Id:      job.UUID,
			Headers: headers,
			Params: &tx.Params{
				From:           job.Transaction.Sender,
				To:             job.Transaction.Recipient,
				Gas:            job.Transaction.Gas,
				GasPrice:       job.Transaction.GasPrice,
				Value:          job.Transaction.Value,
				Nonce:          job.Transaction.Nonce,
				Data:           job.Transaction.Data,
				Raw:            job.Transaction.Raw,
				PrivateFor:     job.Transaction.PrivateFor,
				PrivateFrom:    job.Transaction.PrivateFrom,
				PrivacyGroupId: job.Transaction.PrivacyGroupID,
			},
			ContextLabels: contextLabels,
			JobType:       tx.JobTypeMap[job.Type],
		}},
		InternalLabels: make(map[string]string),
	}

	txEnvelope.SetChainUUID(job.ChainUUID)

	chainID := new(big.Int)
	chainID.SetString(job.InternalData.ChainID, 10)
	txEnvelope.SetChainID(chainID)
	txEnvelope.SetScheduleUUID(job.Schedule.UUID)

	if job.InternalData.OneTimeKey {
		txEnvelope.EnableTxFromOneTimeKey()
	}

	return txEnvelope
}
