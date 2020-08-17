package formatters

import (
	"net/http"
	"strings"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx-scheduler"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

func FormatSendTxRequest(sendTxRequest *types.SendTransactionRequest, idempotencyKey string, defaultRetryInterval time.Duration) *entities.TxRequest {
	return &entities.TxRequest{
		IdempotencyKey: idempotencyKey,
		ChainName:      sendTxRequest.ChainName,
		Labels:         sendTxRequest.Labels,
		Params: &entities.ETHTransactionParams{
			From:            sendTxRequest.Params.From,
			To:              sendTxRequest.Params.To,
			Value:           sendTxRequest.Params.Value,
			GasPrice:        sendTxRequest.Params.GasPrice,
			Gas:             sendTxRequest.Params.Gas,
			MethodSignature: sendTxRequest.Params.MethodSignature,
			Args:            sendTxRequest.Params.Args,
			Protocol:        sendTxRequest.Params.Protocol,
			PrivateFrom:     sendTxRequest.Params.PrivateFrom,
			PrivateFor:      sendTxRequest.Params.PrivateFor,
			PrivacyGroupID:  sendTxRequest.Params.PrivacyGroupID,
		},
		InternalData: formatAnnotations(&sendTxRequest.Params.Annotations, defaultRetryInterval),
	}
}

func FormatDeployContractRequest(deployRequest *types.DeployContractRequest, idempotencyKey string, defaultRetryInterval time.Duration) *entities.TxRequest {
	return &entities.TxRequest{
		IdempotencyKey: idempotencyKey,
		ChainName:      deployRequest.ChainName,
		Labels:         deployRequest.Labels,
		Params: &entities.ETHTransactionParams{
			From:           deployRequest.Params.From,
			Value:          deployRequest.Params.Value,
			GasPrice:       deployRequest.Params.GasPrice,
			Gas:            deployRequest.Params.Gas,
			Args:           deployRequest.Params.Args,
			ContractName:   deployRequest.Params.ContractName,
			ContractTag:    deployRequest.Params.ContractTag,
			Protocol:       deployRequest.Params.Protocol,
			PrivateFrom:    deployRequest.Params.PrivateFrom,
			PrivateFor:     deployRequest.Params.PrivateFor,
			PrivacyGroupID: deployRequest.Params.PrivacyGroupID,
		},
		InternalData: formatAnnotations(&deployRequest.Params.Annotations, defaultRetryInterval),
	}
}

func FormatSendRawRequest(rawTxRequest *types.RawTransactionRequest, idempotencyKey string, defaultRetryInterval time.Duration) *entities.TxRequest {
	// Do not use InternalData directly as we only want to expose the RetryInterval param
	annotations := &types.Annotations{
		RetryPolicy: types.GasPriceRetryParams{
			Interval: rawTxRequest.Params.RetryPolicy.Interval,
		},
	}
	return &entities.TxRequest{
		IdempotencyKey: idempotencyKey,
		ChainName:      rawTxRequest.ChainName,
		Labels:         rawTxRequest.Labels,
		Params: &entities.ETHTransactionParams{
			Raw: rawTxRequest.Params.Raw,
		},
		InternalData: formatAnnotations(annotations, defaultRetryInterval),
	}
}

func FormatTransferRequest(transferRequest *types.TransferRequest, idempotencyKey string, defaultRetryInterval time.Duration) *entities.TxRequest {
	// Do not use InternalData directly as we do not want to expose the OneTimeKey param
	annotations := &types.Annotations{
		Priority:    transferRequest.Params.Priority,
		RetryPolicy: transferRequest.Params.RetryPolicy,
	}
	return &entities.TxRequest{
		IdempotencyKey: idempotencyKey,
		ChainName:      transferRequest.ChainName,
		Labels:         transferRequest.Labels,
		Params: &entities.ETHTransactionParams{
			From:     transferRequest.Params.From,
			To:       transferRequest.Params.To,
			Value:    transferRequest.Params.Value,
			GasPrice: transferRequest.Params.GasPrice,
			Gas:      transferRequest.Params.Gas,
		},
		InternalData: formatAnnotations(annotations, defaultRetryInterval),
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
