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
	mapParams, err := jsonToMap(txRequestModel.Params)
	if err != nil {
		return nil, err
	}

	return &types.TransactionResponse{
		IdempotencyKey: txRequestModel.IdempotencyKey,
		Params:         mapParams,
		Schedule: types.ScheduleResponse{
			UUID:      txRequestModel.Schedule.UUID,
			ChainID:   txRequestModel.Schedule.ChainID,
			CreatedAt: txRequestModel.Schedule.CreatedAt,
		},
		CreatedAt: txRequestModel.CreatedAt,
	}, nil
}

func jsonToMap(jsonStr string) (map[string]interface{}, error) {
	jsonMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &jsonMap)
	if err != nil {
		errMessage := "could not unmarshal JSON string"
		log.WithError(err).Error(errMessage)
		return nil, errors.InvalidFormatError(errMessage)
	}

	return jsonMap, nil
}
