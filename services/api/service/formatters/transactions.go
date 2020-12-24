package formatters

import (
	"net/http"
	"strings"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"

	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
)

func FormatSendTxRequest(sendTxRequest *types.SendTransactionRequest, idempotencyKey string) *entities.TxRequest {
	txRequest := &entities.TxRequest{
		IdempotencyKey: idempotencyKey,
		ChainName:      sendTxRequest.ChainName,
		Labels:         sendTxRequest.Labels,
		Params: &entities.ETHTransactionParams{
			To:              ethcommon.HexToAddress(sendTxRequest.Params.To).Hex(),
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
		InternalData: buildInternalData(
			sendTxRequest.Params.OneTimeKey,
			&sendTxRequest.Params.GasPricePolicy,
		),
	}

	if sendTxRequest.Params.From != "" {
		txRequest.Params.From = ethcommon.HexToAddress(sendTxRequest.Params.From).Hex()
	}

	return txRequest
}

func FormatDeployContractRequest(deployRequest *types.DeployContractRequest, idempotencyKey string) *entities.TxRequest {
	txRequest := &entities.TxRequest{
		IdempotencyKey: idempotencyKey,
		ChainName:      deployRequest.ChainName,
		Labels:         deployRequest.Labels,
		Params: &entities.ETHTransactionParams{
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
		InternalData: buildInternalData(
			deployRequest.Params.OneTimeKey,
			&deployRequest.Params.GasPricePolicy,
		),
	}

	if deployRequest.Params.From != "" {
		txRequest.Params.From = ethcommon.HexToAddress(deployRequest.Params.From).Hex()
	}

	return txRequest
}

func FormatSendRawRequest(rawTxRequest *types.RawTransactionRequest, idempotencyKey string) *entities.TxRequest {
	// Do not use InternalData directly as we only want to expose the RetryInterval param
	gasPricePolicy := &types.GasPriceParams{
		RetryPolicy: types.RetryParams{
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
		InternalData: buildInternalData(false, gasPricePolicy),
	}
}

func FormatTransferRequest(transferRequest *types.TransferRequest, idempotencyKey string) *entities.TxRequest {
	return &entities.TxRequest{
		IdempotencyKey: idempotencyKey,
		ChainName:      transferRequest.ChainName,
		Labels:         transferRequest.Labels,
		Params: &entities.ETHTransactionParams{
			From:     ethcommon.HexToAddress(transferRequest.Params.From).Hex(),
			To:       ethcommon.HexToAddress(transferRequest.Params.To).Hex(),
			Value:    transferRequest.Params.Value,
			GasPrice: transferRequest.Params.GasPrice,
			Gas:      transferRequest.Params.Gas,
		},
		InternalData: buildInternalData(
			false,
			&transferRequest.Params.GasPricePolicy,
		),
	}
}

func FormatTxResponse(txRequest *entities.TxRequest) *types.TransactionResponse {
	scheduleRes := FormatScheduleResponse(txRequest.Schedule)

	return &types.TransactionResponse{
		UUID:           txRequest.Schedule.UUID,
		IdempotencyKey: txRequest.IdempotencyKey,
		ChainName:      txRequest.ChainName,
		Params:         txRequest.Params,
		Jobs:           scheduleRes.Jobs,
		CreatedAt:      txRequest.CreatedAt,
	}
}

func FormatTransactionsFilterRequest(req *http.Request) (*entities.TransactionRequestFilters, error) {
	filters := &entities.TransactionRequestFilters{}

	qIdempotencyKeys := req.URL.Query().Get("idempotency_keys")
	if qIdempotencyKeys != "" {
		filters.IdempotencyKeys = strings.Split(qIdempotencyKeys, ",")
	}

	if err := utils.GetValidator().Struct(filters); err != nil {
		return nil, err
	}

	return filters, nil
}

func buildInternalData(oneTimeKey bool, gasPricePolicy *types.GasPriceParams) *entities.InternalData {
	internalData := &entities.InternalData{
		OneTimeKey:        oneTimeKey,
		Priority:          gasPricePolicy.Priority,
		GasPriceIncrement: gasPricePolicy.RetryPolicy.Increment,
		GasPriceLimit:     gasPricePolicy.RetryPolicy.Limit,
	}

	if gasPricePolicy.RetryPolicy.Interval != "" {
		// we can skip the error check as at this point we know the interval is a duration as it already passed validation
		internalData.RetryInterval, _ = time.ParseDuration(gasPricePolicy.RetryPolicy.Interval)
	}

	if internalData.Priority == "" {
		internalData.Priority = utils.PriorityMedium
	}

	return internalData
}
