package parsers

import (
	"fmt"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/services/api/store/models"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
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

func NewJobEntitiesFromTxRequest(txRequest *entities.TxRequest, chainUUID, txData string) ([]*entities.Job, error) {
	var jobs []*entities.Job
	switch {
	case txRequest.Params.Protocol == entities.EEAChainType:
		privTxJob := newJobEntityFromTxRequest(txRequest, newEthTransactionFromParams(txRequest.Params, txData, entities.LegacyTxType), entities.EEAPrivateTransaction, chainUUID)
		markingTxJob := newJobEntityFromTxRequest(txRequest, &entities.ETHTransaction{}, entities.EEAMarkingTransaction, chainUUID)
		markingTxJob.InternalData.OneTimeKey = true
		jobs = append(jobs, privTxJob, markingTxJob)
	case txRequest.Params.Protocol == entities.TesseraChainType:
		privTxJob := newJobEntityFromTxRequest(txRequest, newEthTransactionFromParams(txRequest.Params, txData, entities.LegacyTxType),
			entities.TesseraPrivateTransaction, chainUUID)

		markingTx := &entities.ETHTransaction{
			From:         "",
			PrivateFor:   txRequest.Params.PrivateFor,
			MandatoryFor: txRequest.Params.MandatoryFor,
			PrivacyFlag:  txRequest.Params.PrivacyFlag,
		}
		if txRequest.Params.From != nil {
			markingTx.From = txRequest.Params.From.Hex()
		}
		markingTxJob := newJobEntityFromTxRequest(txRequest, markingTx, entities.TesseraMarkingTransaction, chainUUID)
		jobs = append(jobs, privTxJob, markingTxJob)
	case txRequest.Params.Raw != "":
		rawTx, err := newTransactionFromRaw(txRequest.Params.Raw)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, newJobEntityFromTxRequest(txRequest, rawTx, entities.EthereumRawTransaction, chainUUID))
	default:
		tx := newEthTransactionFromParams(txRequest.Params, txData, entities.TransactionType(txRequest.Params.TransactionType))
		jobs = append(jobs, newJobEntityFromTxRequest(txRequest, tx, entities.EthereumTransaction, chainUUID))
	}

	return jobs, nil
}

func newEthTransactionFromParams(params *entities.ETHTransactionParams, txData string, txType entities.TransactionType) *entities.ETHTransaction {
	tx := &entities.ETHTransaction{
		From:            "",
		To:              "",
		Nonce:           params.Nonce,
		Value:           params.Value,
		GasPrice:        params.GasPrice,
		Gas:             params.Gas,
		GasFeeCap:       params.GasFeeCap,
		GasTipCap:       params.GasTipCap,
		AccessList:      params.AccessList,
		TransactionType: txType,
		Raw:             params.Raw,
		Data:            txData,
		PrivateFrom:     params.PrivateFrom,
		PrivateFor:      params.PrivateFor,
		MandatoryFor:    params.MandatoryFor,
		PrivacyFlag:     params.PrivacyFlag,
		PrivacyGroupID:  params.PrivacyGroupID,
	}
	if params.From != nil {
		tx.From = params.From.Hex()
	}
	if params.To != nil {
		tx.To = params.To.Hex()
	}
	return tx
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
		OwnerID:      txRequest.Schedule.OwnerID,
	}
}

func newTransactionFromRaw(raw string) (*entities.ETHTransaction, error) {
	tx := &types.Transaction{}

	rawb, err := hexutil.Decode(raw)
	if err != nil {
		return nil, err
	}

	err = tx.UnmarshalBinary(rawb)
	if err != nil {
		return nil, errors.InvalidParameterError(err.Error())
	}

	from, err := types.Sender(types.NewEIP155Signer(tx.ChainId()), tx)
	if err != nil {
		return nil, errors.InvalidParameterError(err.Error())
	}

	jobTx := &entities.ETHTransaction{
		From:     from.String(),
		Data:     hexutil.Encode(tx.Data()),
		Gas:      fmt.Sprintf("%d", tx.Gas()),
		GasPrice: fmt.Sprintf("%d", tx.GasPrice()),
		Value:    tx.Value().String(),
		Nonce:    fmt.Sprintf("%d", tx.Nonce()),
		Hash:     tx.Hash().String(),
		Raw:      raw,
	}

	// If not contract creation
	if tx.To() != nil {
		jobTx.To = tx.To().String()
	}

	return jobTx, nil
}
