package formatters

import (
	"net/http"
	"strings"
	"time"

	"github.com/consensys/orchestrate/pkg/errors"

	types "github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/pkg/types/entities"

	"github.com/consensys/orchestrate/pkg/utils"
)

func FormatSendTxRequest(sendTxRequest *types.SendTransactionRequest, idempotencyKey string) *entities.TxRequest {
	if sendTxRequest.Params.ContractTag == "" {
		sendTxRequest.Params.ContractTag = entities.DefaultTagValue
	}

	txRequest := &entities.TxRequest{
		IdempotencyKey: idempotencyKey,
		ChainName:      sendTxRequest.ChainName,
		Labels:         sendTxRequest.Labels,
		Params: &entities.ETHTransactionParams{
			From:            sendTxRequest.Params.From,
			To:              sendTxRequest.Params.To,
			Nonce:           sendTxRequest.Params.Nonce,
			Value:           sendTxRequest.Params.Value,
			GasPrice:        sendTxRequest.Params.GasPrice,
			Gas:             sendTxRequest.Params.Gas,
			GasFeeCap:       sendTxRequest.Params.GasFeeCap,
			GasTipCap:       sendTxRequest.Params.GasTipCap,
			AccessList:      sendTxRequest.Params.AccessList,
			TransactionType: sendTxRequest.Params.TransactionType,
			MethodSignature: sendTxRequest.Params.MethodSignature,
			Args:            sendTxRequest.Params.Args,
			Protocol:        sendTxRequest.Params.Protocol,
			PrivateFrom:     sendTxRequest.Params.PrivateFrom,
			PrivateFor:      sendTxRequest.Params.PrivateFor,
			MandatoryFor:    sendTxRequest.Params.MandatoryFor,
			PrivacyFlag:     sendTxRequest.Params.PrivacyFlag,
			PrivacyGroupID:  sendTxRequest.Params.PrivacyGroupID,
			ContractTag:     sendTxRequest.Params.ContractTag,
			ContractName:    sendTxRequest.Params.ContractName,
		},
		InternalData: buildInternalData(
			sendTxRequest.Params.OneTimeKey,
			&sendTxRequest.Params.GasPricePolicy,
		),
	}

	return txRequest
}

func FormatDeployContractRequest(deployRequest *types.DeployContractRequest, idempotencyKey string) *entities.TxRequest {
	if deployRequest.Params.ContractTag == "" {
		deployRequest.Params.ContractTag = entities.DefaultTagValue
	}

	txRequest := &entities.TxRequest{
		IdempotencyKey: idempotencyKey,
		ChainName:      deployRequest.ChainName,
		Labels:         deployRequest.Labels,
		Params: &entities.ETHTransactionParams{
			From:            deployRequest.Params.From,
			Nonce:           deployRequest.Params.Nonce,
			Value:           deployRequest.Params.Value,
			GasPrice:        deployRequest.Params.GasPrice,
			Gas:             deployRequest.Params.Gas,
			GasFeeCap:       deployRequest.Params.GasFeeCap,
			GasTipCap:       deployRequest.Params.GasTipCap,
			AccessList:      deployRequest.Params.AccessList,
			TransactionType: deployRequest.Params.TransactionType,
			Args:            deployRequest.Params.Args,
			ContractName:    deployRequest.Params.ContractName,
			ContractTag:     deployRequest.Params.ContractTag,
			Protocol:        deployRequest.Params.Protocol,
			PrivateFrom:     deployRequest.Params.PrivateFrom,
			PrivateFor:      deployRequest.Params.PrivateFor,
			MandatoryFor:    deployRequest.Params.MandatoryFor,
			PrivacyFlag:     entities.PrivacyFlag(deployRequest.Params.PrivacyFlag),
			PrivacyGroupID:  deployRequest.Params.PrivacyGroupID,
		},
		InternalData: buildInternalData(
			deployRequest.Params.OneTimeKey,
			&deployRequest.Params.GasPricePolicy,
		),
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
			From:            &transferRequest.Params.From,
			To:              &transferRequest.Params.To,
			Nonce:           transferRequest.Params.Nonce,
			GasFeeCap:       transferRequest.Params.GasFeeCap,
			GasTipCap:       transferRequest.Params.GasTipCap,
			AccessList:      transferRequest.Params.AccessList,
			TransactionType: transferRequest.Params.TransactionType,
			Value:           transferRequest.Params.Value,
			GasPrice:        transferRequest.Params.GasPrice,
			Gas:             transferRequest.Params.Gas,
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

	pagination, err := utils.FilterIntegerValueWithKey(req)
	if err != nil {
		return filters, err
	}

	filters.Pagination = *pagination

	if err := utils.GetValidator().Struct(filters); err != nil {
		return nil, errors.InvalidFormatError(err.Error())
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
