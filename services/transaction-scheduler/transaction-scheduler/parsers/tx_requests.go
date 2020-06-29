package parsers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
)

func NewTxRequestModelFromEntities(txRequest *entities.TxRequest, requestHash string) *models.TransactionRequest {
	return &models.TransactionRequest{
		UUID:           txRequest.UUID,
		IdempotencyKey: txRequest.IdempotencyKey,
		RequestHash:    requestHash,
		Params:         txRequest.Params,
		CreatedAt:      txRequest.CreatedAt,
	}
}

func NewJobEntityFromTxRequest(txRequest *entities.TxRequest, jobType, chainUUID string) *types.Job {
	return &types.Job{
		ScheduleUUID: txRequest.Schedule.UUID,
		ChainUUID:    chainUUID,
		Type:         jobType,
		Labels:       txRequest.Labels,
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
}
