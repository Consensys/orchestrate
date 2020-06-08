package parsers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
)

func NewScheduleEntityFromModels(scheduleModel *models.Schedule) *entities.Schedule {
	schedule := &entities.Schedule{
		UUID:      scheduleModel.UUID,
		CreatedAt: scheduleModel.CreatedAt,
	}

	for _, job := range scheduleModel.Jobs {
		schedule.Jobs = append(schedule.Jobs, NewJobEntityFromModels(job))
	}

	return schedule
}

func NewScheduleModelFromEntities(schedule *entities.Schedule, tenantID string) *models.Schedule {
	scheduleModel := &models.Schedule{
		UUID:     schedule.UUID,
		TenantID: tenantID,
	}

	if schedule.TxRequest != nil {
		scheduleModel.TransactionRequest = &models.TransactionRequest{
			IdempotencyKey: schedule.TxRequest.IdempotencyKey,
		}
	}

	for _, job := range schedule.Jobs {
		scheduleModel.Jobs = append(scheduleModel.Jobs, NewJobModelFromEntities(job, &scheduleModel.ID))
	}

	return scheduleModel
}
