package parsers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
)

func NewLogEntityFromModels(logModel *models.Log) *types.Log {
	return &types.Log{
		Status:    logModel.Status,
		Message:   logModel.Message,
		CreatedAt: logModel.CreatedAt,
	}
}

func NewLogModelFromEntity(log *types.Log) *models.Log {
	return &models.Log{
		Status:    log.Status,
		Message:   log.Message,
		CreatedAt: log.CreatedAt,
	}
}
