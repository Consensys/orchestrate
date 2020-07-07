package parsers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
)

func NewTxRequestModelFromEntities(txRequest *entities.TxRequest, requestHash string, scheduleID int) *models.TransactionRequest {
	return &models.TransactionRequest{
		UUID:           txRequest.UUID,
		IdempotencyKey: txRequest.IdempotencyKey,
		RequestHash:    requestHash,
		Params:         txRequest.Params,
		ScheduleID:     &scheduleID,
		CreatedAt:      txRequest.CreatedAt,
	}
}

func NewJobEntityFromTxRequest(txRequest *entities.TxRequest, jobType, chainUUID string) *types.Job {
	job := &types.Job{
		ScheduleUUID: txRequest.Schedule.UUID,
		ChainUUID:    chainUUID,
		Type:         jobType,
		Labels:       txRequest.Labels,
		Annotations:  &types.Annotations{},
		Transaction: &types.ETHTransaction{
			From:           txRequest.Params.From,
			To:             txRequest.Params.To,
			Nonce:          txRequest.Params.Nonce,
			Value:          txRequest.Params.Value,
			GasPrice:       txRequest.Params.GasPrice,
			Gas:            txRequest.Params.Gas,
			Raw:            txRequest.Params.Raw,
			PrivateFrom:    txRequest.Params.PrivateFrom,
			PrivateFor:     txRequest.Params.PrivateFor,
			PrivacyGroupID: txRequest.Params.PrivacyGroupID,
		},
	}

	if txRequest.Annotations != nil {
		job.Annotations = txRequest.Annotations
	}

	return job
}
