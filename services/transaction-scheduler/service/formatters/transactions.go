package formatters

import (
	types2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
)

func FormatSendTxRequest(txRequest *types.SendTransactionRequest) *entities.TxRequest {
	return &entities.TxRequest{
		IdempotencyKey: txRequest.IdempotencyKey,
		Labels:         txRequest.Labels,
		Params: &types2.ETHTransactionParams{
			From:            txRequest.Params.From,
			To:              txRequest.Params.To,
			Value:           txRequest.Params.Value,
			GasPrice:        txRequest.Params.GasPrice,
			GasLimit:        txRequest.Params.Gas,
			MethodSignature: txRequest.Params.MethodSignature,
			Args:            txRequest.Params.Args,
		},
	}
}

func FormatTxResponse(txRequest *entities.TxRequest) *types.TransactionResponse {
	return &types.TransactionResponse{
		IdempotencyKey: txRequest.IdempotencyKey,
		Params:         txRequest.Params,
		Schedule:       FormatScheduleResponse(txRequest.Schedule),
		CreatedAt:      txRequest.CreatedAt,
	}
}
