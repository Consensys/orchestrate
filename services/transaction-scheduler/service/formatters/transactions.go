package formatters

import (
	pkgtypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
)

func FormatSendTxRequest(txRequest *types.SendTransactionRequest) *entities.TxRequest {
	return &entities.TxRequest{
		IdempotencyKey: txRequest.IdempotencyKey,
		Labels:         txRequest.Labels,
		Params: &pkgtypes.ETHTransactionParams{
			From:                     txRequest.Params.From,
			To:                       txRequest.Params.To,
			Value:                    txRequest.Params.Value,
			GasPrice:                 txRequest.Params.GasPrice,
			GasLimit:                 txRequest.Params.Gas,
			MethodSignature:          txRequest.Params.MethodSignature,
			Args:                     txRequest.Params.Args,
			PrivateTransactionParams: txRequest.Params.PrivateTransactionParams,
		},
	}
}

func FormatDeployContractRequest(txRequest *types.DeployContractRequest) *entities.TxRequest {
	return &entities.TxRequest{
		IdempotencyKey: txRequest.IdempotencyKey,
		Labels:         txRequest.Labels,
		Params: &pkgtypes.ETHTransactionParams{
			From:                     txRequest.Params.From,
			Value:                    txRequest.Params.Value,
			GasPrice:                 txRequest.Params.GasPrice,
			GasLimit:                 txRequest.Params.Gas,
			Args:                     txRequest.Params.Args,
			ContractName:             txRequest.Params.ContractName,
			ContractTag:              txRequest.Params.ContractTag,
			PrivateTransactionParams: txRequest.Params.PrivateTransactionParams,
		},
	}
}

func FormatSendRawRequest(txRequest *types.RawTransactionRequest) *entities.TxRequest {
	return &entities.TxRequest{
		IdempotencyKey: txRequest.IdempotencyKey,
		Labels:         txRequest.Labels,
		Params: &pkgtypes.ETHTransactionParams{
			Raw: txRequest.Params.Raw,
		},
	}
}

func FormatTxResponse(txRequest *entities.TxRequest, chainName string) *types.TransactionResponse {
	return &types.TransactionResponse{
		IdempotencyKey: txRequest.IdempotencyKey,
		Params:         txRequest.Params,
		ChainName:      chainName,
		Schedule:       FormatScheduleResponse(txRequest.Schedule),
		CreatedAt:      txRequest.CreatedAt,
	}
}
