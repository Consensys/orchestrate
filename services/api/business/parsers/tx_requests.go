package parsers

import (
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/ConsenSys/orchestrate/services/api/store/models"
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
	case txRequest.Params.Protocol == entities.OrionChainType:
		privTxJob := newJobEntityFromTxRequest(txRequest, newEthTransactionFromParams(txRequest.Params, txData), entities.OrionEEATransaction, chainUUID)
		markingTxJob := newJobEntityFromTxRequest(txRequest, &entities.ETHTransaction{}, entities.OrionMarkingTransaction, chainUUID)
		markingTxJob.InternalData.OneTimeKey = true
		jobs = append(jobs, privTxJob, markingTxJob)
	case txRequest.Params.Protocol == entities.TesseraChainType:
		privTxJob := newJobEntityFromTxRequest(txRequest, newEthTransactionFromParams(txRequest.Params, txData),
			entities.TesseraPrivateTransaction, chainUUID)
		markingTxJob := newJobEntityFromTxRequest(txRequest, &entities.ETHTransaction{From: txRequest.Params.From,
			PrivateFor: txRequest.Params.PrivateFor}, entities.TesseraMarkingTransaction, chainUUID)
		jobs = append(jobs, privTxJob, markingTxJob)
	case txRequest.Params.Raw != "":
		jobs = append(jobs, newJobEntityFromTxRequest(txRequest, newEthTransactionFromParams(txRequest.Params, txData),
			entities.EthereumRawTransaction, chainUUID))
	default:
		jobs = append(jobs, newJobEntityFromTxRequest(txRequest, newEthTransactionFromParams(txRequest.Params, txData),
			entities.EthereumTransaction, chainUUID))
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

func newJobEntityFromTxRequest(txRequest *entities.TxRequest, ethTx *entities.ETHTransaction, jobType entities.JobType, chainUUID string) *entities.Job {
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
