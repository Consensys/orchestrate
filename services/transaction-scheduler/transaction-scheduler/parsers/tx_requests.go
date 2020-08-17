package parsers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
)

func NewTxRequestModelFromEntities(txRequest *entities.TxRequest, requestHash string, scheduleID int) *models.TransactionRequest {
	return &models.TransactionRequest{
		UUID:           txRequest.UUID,
		IdempotencyKey: txRequest.IdempotencyKey,
		ChainName:      txRequest.ChainName,
		RequestHash:    requestHash,
		Params:         txRequest.Params,
		ScheduleID:     &scheduleID,
		CreatedAt:      txRequest.CreatedAt,
	}
}

func NewJobEntityFromTxRequest(txRequest *entities.TxRequest, jobType, chainUUID string) *entities.Job {
	return &entities.Job{
		ScheduleUUID: txRequest.Schedule.UUID,
		ChainUUID:    chainUUID,
		Type:         jobType,
		Labels:       txRequest.Labels,
		InternalData: txRequest.InternalData,
		Transaction: &entities.ETHTransaction{
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
}
