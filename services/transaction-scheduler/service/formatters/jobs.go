package formatters

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
)

func FormatJobResponse(job *entities.Job) *types.JobResponse {
	jobResponse := &types.JobResponse{
		UUID: job.UUID,
		Transaction: types.ETHTransaction{
			Hash:           job.Transaction.Hash,
			From:           job.Transaction.From,
			To:             job.Transaction.To,
			Nonce:          job.Transaction.Nonce,
			Value:          job.Transaction.Value,
			GasPrice:       job.Transaction.GasPrice,
			GasLimit:       job.Transaction.GasLimit,
			Data:           job.Transaction.Data,
			Raw:            job.Transaction.Raw,
			PrivateFrom:    job.Transaction.PrivateFrom,
			PrivateFor:     job.Transaction.PrivateFor,
			PrivacyGroupID: job.Transaction.PrivacyGroupID,
		},
		Status:    job.Status,
		CreatedAt: job.CreatedAt,
	}

	return jobResponse
}

func FormatJobCreateRequest(request *types.CreateJobRequest) *entities.Job {
	job := &entities.Job{
		Type:         request.Type,
		Labels:       request.Labels,
		ScheduleUUID: request.ScheduleUUID,
		Transaction:  formatTxRequest(&request.Transaction),
	}

	return job
}

func FormatJobUpdateRequest(request *types.UpdateJobRequest) *entities.Job {
	job := &entities.Job{
		Labels:      request.Labels,
		Transaction: formatTxRequest(&request.Transaction),
	}

	return job
}

func formatTxRequest(tx *types.ETHTransaction) *entities.Transaction {
	return &entities.Transaction{
		Hash:           tx.Hash,
		From:           tx.From,
		To:             tx.To,
		Nonce:          tx.Nonce,
		Value:          tx.Value,
		GasPrice:       tx.GasPrice,
		GasLimit:       tx.GasLimit,
		Data:           tx.Data,
		PrivateFrom:    tx.PrivateFrom,
		PrivateFor:     tx.PrivateFor,
		PrivacyGroupID: tx.PrivacyGroupID,
		Raw:            tx.Raw,
	}
}
