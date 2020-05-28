package formatters

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
)

func FormatSendTxRequest(txRequest *types.SendTransactionRequest, chainUUID string) *entities.TxRequest {
	return &entities.TxRequest{
		IdempotencyKey: txRequest.IdempotencyKey,
		Labels:         txRequest.Labels,
		Schedule: &entities.Schedule{
			ChainUUID: chainUUID,
			Jobs: []*entities.Job{{
				Type: tx.JobEthereumTransaction,
			}},
		},
		Params: &entities.TxRequestParams{
			From:            &txRequest.Params.From,
			To:              &txRequest.Params.To,
			Value:           &txRequest.Params.Value,
			GasPrice:        &txRequest.Params.GasPrice,
			MethodSignature: &txRequest.Params.MethodSignature,
			Args:            txRequest.Params.Args,
		},
	}
}

func FormatTxResponse(txRequest *entities.TxRequest) (*types.TransactionResponse, error) {
	jsonStr, err := json.Marshal(txRequest.Params)
	if err != nil {
		return nil, err
	}

	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(jsonStr, &jsonMap)
	if err != nil {
		return nil, err
	}

	return &types.TransactionResponse{
		IdempotencyKey: txRequest.IdempotencyKey,
		Params:         jsonMap,
		Schedule:       FormatScheduleResponse(txRequest.Schedule),
		CreatedAt:      txRequest.CreatedAt,
	}, nil
}
