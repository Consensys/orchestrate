package parsers

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
)

func NewTxRequestModelFromEntities(txRequest *entities.TxRequest, requestHash, tenantID string) (*models.TransactionRequest, error) {
	scheduleModel := &models.Schedule{
		TenantID:  tenantID,
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
	}

	return txRequestModel, nil
}

func NewJobEntityFromSendTxRequest(txRequest *entities.TxRequest) (*entities.Job, error) {
	txEntity := &entities.Transaction{}

	if txRequest.Params.From == nil {
		return nil, errors.InvalidArgError("missing required param '%s'", "From")
	}
	if txRequest.Params.To == nil {
		return nil, errors.InvalidArgError("missing required param '%s'", "To")
	}

	txEntity.From = *txRequest.Params.From
	txEntity.To = *txRequest.Params.To
	if txRequest.Params.Value != nil {
		txEntity.Value = *txRequest.Params.Value
	}
	if txRequest.Params.GasPrice != nil {
		txEntity.GasPrice = *txRequest.Params.GasPrice
	}
	if txRequest.Params.GasLimit != nil {
		txEntity.GasLimit = *txRequest.Params.GasLimit
	}

	if txRequest.Params.MethodSignature == nil {
		return nil, errors.InvalidArgError("missing required param '%s'", "MethodSignature")
	}

	crafter := abi.BaseCrafter{}
	txDataBytes, err := crafter.CraftCall(*txRequest.Params.MethodSignature, txRequest.Params.Args...)
	if err != nil {
		return nil, err
	}

	txEntity.Data = hexutil.Encode(txDataBytes)

	return &entities.Job{
		ScheduleUUID: txRequest.Schedule.UUID,
		Type:         entities.JobConstantinopleTransaction,
		Labels:       txRequest.Labels,
		Transaction:  txEntity,
	}, nil
}
