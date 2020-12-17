package parsers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models"
)

func NewTxRequestModelFromEntities(txRequest *entities.TxRequest, requestHash string, scheduleID int) *models.TransactionRequest {
	return &models.TransactionRequest{
		IdempotencyKey: txRequest.IdempotencyKey,
		ChainName:      txRequest.ChainName,
		RequestHash:    requestHash,
		Params:         txRequest.Params,
		ScheduleID:     &scheduleID,
		CreatedAt:      txRequest.CreatedAt,
	}
}

func NewJobEntitiesFromTxRequest(txRequest *entities.TxRequest, chainUUID, txData string) []*entities.Job {
	var jobs []*entities.Job
	switch {
	case txRequest.Params.Protocol == utils.OrionChainType:
		privTxJob := newJobEntityFromTxRequest(txRequest, newEthTransactionFromParams(txRequest.Params, txData), utils.OrionEEATransaction, chainUUID)
		markingTxJob := newJobEntityFromTxRequest(txRequest, &entities.ETHTransaction{}, utils.OrionMarkingTransaction, chainUUID)
		markingTxJob.InternalData.OneTimeKey = true
		jobs = append(jobs, privTxJob, markingTxJob)
	case txRequest.Params.Protocol == utils.TesseraChainType:
		privTxJob := newJobEntityFromTxRequest(txRequest, newEthTransactionFromParams(txRequest.Params, txData), utils.TesseraPrivateTransaction, chainUUID)
		markingTxJob := newJobEntityFromTxRequest(txRequest, &entities.ETHTransaction{From: txRequest.Params.From, PrivateFor: txRequest.Params.PrivateFor}, utils.TesseraMarkingTransaction, chainUUID)
		jobs = append(jobs, privTxJob, markingTxJob)
	case txRequest.Params.Raw != "":
		jobs = append(jobs, newJobEntityFromTxRequest(txRequest, newEthTransactionFromParams(txRequest.Params, txData),
			utils.EthereumRawTransaction, chainUUID))
	default:
		jobs = append(jobs, newJobEntityFromTxRequest(txRequest, newEthTransactionFromParams(txRequest.Params, txData),
			utils.EthereumTransaction, chainUUID))
	}

	return jobs
}

func newEthTransactionFromParams(params *entities.ETHTransactionParams, txData string) *entities.ETHTransaction {
	return &entities.ETHTransaction{
		From:           params.From,
		To:             params.To,
		Nonce:          params.Nonce,
		Value:          params.Value,
		GasPrice:       params.GasPrice,
		Gas:            params.Gas,
		Raw:            params.Raw,
		Data:           txData,
		PrivateFrom:    params.PrivateFrom,
		PrivateFor:     params.PrivateFor,
		PrivacyGroupID: params.PrivacyGroupID,
	}
}

func newJobEntityFromTxRequest(txRequest *entities.TxRequest, ethTx *entities.ETHTransaction, jobType, chainUUID string) *entities.Job {
	internalData := *txRequest.InternalData
	return &entities.Job{
		ScheduleUUID: txRequest.Schedule.UUID,
		ChainUUID:    chainUUID,
		Type:         jobType,
		Labels:       txRequest.Labels,
		InternalData: &internalData,
		Transaction:  ethTx,
		TenantID:     txRequest.Schedule.TenantID,
	}
}
