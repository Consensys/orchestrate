package parsers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/store/models"
)

func NewLogEntityFromModels(logModel *models.Log) *entities.Log {
	return &entities.Log{
		Status:    logModel.Status,
		Message:   logModel.Message,
		CreatedAt: logModel.CreatedAt,
	}
}

func NewLogModelFromEntity(log *entities.Log) *models.Log {
	return &models.Log{
		Status:    log.Status,
		Message:   log.Message,
		CreatedAt: log.CreatedAt,
	}
}
