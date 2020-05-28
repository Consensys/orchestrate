package parsers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
)

func NewTxRequestModelFromEntities(txRequest *entities.TxRequest, requestHash, tenantID string) (*models.TransactionRequest, error) {
	scheduleModel := &models.Schedule{
		TenantID:  tenantID,
		UUID:      txRequest.Schedule.UUID,
		ChainUUID: txRequest.Schedule.ChainUUID,
	}

	jsonParams, err := json.Marshal(txRequest.Params)
	if err != nil {
		return nil, err
	}

	txRequestModel := &models.TransactionRequest{
		IdempotencyKey: txRequest.IdempotencyKey,
		RequestHash:    requestHash,
		Params:         string(jsonParams),
		Schedules:      []*models.Schedule{scheduleModel},
		CreatedAt:      txRequest.CreatedAt,
	}

	return txRequestModel, nil
}

func NewJobEntityFromTxRequest(txRequest *entities.TxRequest, jobType string) *entities.Job {
	txEntity := &entities.ETHTransaction{
		From:           txRequest.Params.From,
		To:             txRequest.Params.To,
		Nonce:          txRequest.Params.Nonce,
		Value:          txRequest.Params.Value,
		GasPrice:       txRequest.Params.GasPrice,
		GasLimit:       txRequest.Params.GasLimit,
		Raw:            txRequest.Params.Raw,
		PrivateFrom:    txRequest.Params.PrivateFrom,
		PrivateFor:     txRequest.Params.PrivateFor,
		PrivacyGroupID: txRequest.Params.PrivacyGroupID,
	}

	return &entities.Job{
		ScheduleUUID: txRequest.Schedule.UUID,
		Type:         jobType,
		Labels:       txRequest.Labels,
		Transaction:  txEntity,
	}
}
