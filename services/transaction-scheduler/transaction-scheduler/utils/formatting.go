package utils

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
)

func ObjectToJSON(obj interface{}) (string, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		errMessage := "could not marshal object"
		log.WithError(err).Error(errMessage)
		return "", errors.InvalidParameterError(errMessage)
	}

	return string(b), nil
}

func FormatTxResponse(txRequestModel *models.TransactionRequest) (*types.TransactionResponse, error) {
	mapParams, err := JSONToMap(txRequestModel.Params)
	if err != nil {
		return nil, err
	}

	return &types.TransactionResponse{
		IdempotencyKey: txRequestModel.IdempotencyKey,
		Params:         mapParams,
		Schedule:       FormatScheduleResponse(txRequestModel.Schedule),
		CreatedAt:      txRequestModel.CreatedAt,
	}, nil
}

func FormatScheduleResponse(scheduleModel *models.Schedule) *types.ScheduleResponse {
	scheduleResponse := &types.ScheduleResponse{
		UUID:      scheduleModel.UUID,
		ChainUUID: scheduleModel.ChainUUID,
		CreatedAt: scheduleModel.CreatedAt,
		Jobs:      []*types.JobResponse{},
	}

	for _, job := range scheduleModel.Jobs {
		scheduleResponse.Jobs = append(scheduleResponse.Jobs, FormatJobResponse(job))
	}

	return scheduleResponse
}

func FormatJobResponse(jobModel *models.Job) *types.JobResponse {
	jobResponse := &types.JobResponse{
		UUID: jobModel.UUID,
		Transaction: types.ETHTransaction{
			Hash:           jobModel.Transaction.Hash,
			From:           jobModel.Transaction.Sender,
			To:             jobModel.Transaction.Recipient,
			Nonce:          jobModel.Transaction.Nonce,
			Value:          jobModel.Transaction.Value,
			GasPrice:       jobModel.Transaction.GasPrice,
			GasLimit:       jobModel.Transaction.GasLimit,
			Data:           jobModel.Transaction.Data,
			Raw:            jobModel.Transaction.Raw,
			PrivateFrom:    jobModel.Transaction.PrivateFrom,
			PrivateFor:     jobModel.Transaction.PrivateFor,
			PrivacyGroupID: jobModel.Transaction.PrivacyGroupID,
		},
		Status:    jobModel.GetStatus(),
		CreatedAt: jobModel.CreatedAt,
	}

	return jobResponse
}

func JSONToMap(jsonStr string) (map[string]interface{}, error) {
	jsonMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &jsonMap)
	if err != nil {
		errMessage := "could not unmarshal JSON string"
		log.WithError(err).Error(errMessage)
		return nil, errors.InvalidFormatError(errMessage)
	}

	return jsonMap, nil
}
