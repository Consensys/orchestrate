package formatters

import (
	"net/http"
	"strings"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
)

func FormatSendTxRequest(txRequest *types.SendTransactionRequest, idempotencyKey string) *entities.TxRequest {
	return &entities.TxRequest{
		IdempotencyKey: idempotencyKey,
		ChainName:      txRequest.ChainName,
		Labels:         txRequest.Labels,
		Params: &types.ETHTransactionParams{
			From:                     txRequest.Params.From,
			To:                       txRequest.Params.To,
			Value:                    txRequest.Params.Value,
			GasPrice:                 txRequest.Params.GasPrice,
			Gas:                      txRequest.Params.Gas,
			MethodSignature:          txRequest.Params.MethodSignature,
			Args:                     txRequest.Params.Args,
			PrivateTransactionParams: txRequest.Params.PrivateTransactionParams,
		},
		Annotations: &types.Annotations{
			OneTimeKey: txRequest.Params.OneTimeKey,
			Priority:   txRequest.Params.Priority,
			Retry:      txRequest.Params.Retry,
		},
	}
}

func FormatDeployContractRequest(txRequest *types.DeployContractRequest, idempotencyKey string) *entities.TxRequest {
	return &entities.TxRequest{
		IdempotencyKey: idempotencyKey,
		ChainName:      txRequest.ChainName,
		Labels:         txRequest.Labels,
		Params: &types.ETHTransactionParams{
			From:                     txRequest.Params.From,
			Value:                    txRequest.Params.Value,
			GasPrice:                 txRequest.Params.GasPrice,
			Gas:                      txRequest.Params.Gas,
			Args:                     txRequest.Params.Args,
			ContractName:             txRequest.Params.ContractName,
			ContractTag:              txRequest.Params.ContractTag,
			PrivateTransactionParams: txRequest.Params.PrivateTransactionParams,
		},
		Annotations: &types.Annotations{
			OneTimeKey: txRequest.Params.OneTimeKey,
			Priority:   txRequest.Params.Priority,
			Retry:      txRequest.Params.Retry,
		},
	}
}

func FormatSendRawRequest(txRequest *types.RawTransactionRequest, idempotencyKey string) *entities.TxRequest {
	var retry *types.GasPriceRetryParams
	if r := txRequest.Params.Retry; r != nil {
		retry = &types.GasPriceRetryParams{BaseRetryParams: types.BaseRetryParams{Interval: r.Interval}}
	}
	return &entities.TxRequest{
		IdempotencyKey: idempotencyKey,
		ChainName:      txRequest.ChainName,
		Labels:         txRequest.Labels,
		Params: &types.ETHTransactionParams{
			Raw: txRequest.Params.Raw,
		},
		Annotations: &types.Annotations{
			Retry: retry,
		},
	}
}

func FormatSendTransferRequest(txRequest *types.TransferRequest, idempotencyKey string) *entities.TxRequest {
	return &entities.TxRequest{
		IdempotencyKey: idempotencyKey,
		ChainName:      txRequest.ChainName,
		Labels:         txRequest.Labels,
		Params: &types.ETHTransactionParams{
			From:     txRequest.Params.From,
			To:       txRequest.Params.To,
			Value:    txRequest.Params.Value,
			GasPrice: txRequest.Params.GasPrice,
			Gas:      txRequest.Params.Gas,
		},
		Annotations: &types.Annotations{
			Priority: txRequest.Params.Priority,
			Retry:    txRequest.Params.Retry,
		},
	}
}

func FormatTxResponse(txRequest *entities.TxRequest) *types.TransactionResponse {
	return &types.TransactionResponse{
		UUID:           txRequest.UUID,
		IdempotencyKey: txRequest.IdempotencyKey,
		ChainName:      txRequest.ChainName,
		Params:         txRequest.Params,
		Schedule:       FormatScheduleResponse(txRequest.Schedule),
		CreatedAt:      txRequest.CreatedAt,
	}
}

func FormatTransactionsFilterRequest(req *http.Request) (*entities.TransactionFilters, error) {
	filters := &entities.TransactionFilters{}

	qIdempotencyKeys := req.URL.Query().Get("idempotency_keys")
	if qIdempotencyKeys != "" {
		filters.IdempotencyKeys = strings.Split(qIdempotencyKeys, ",")
	}

	if err := utils.GetValidator().Struct(filters); err != nil {
		return nil, err
	}

	return filters, nil
}
