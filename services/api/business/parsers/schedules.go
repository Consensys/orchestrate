package parsers

import (
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/services/api/store/models"
)

func NewScheduleEntityFromModels(scheduleModel *models.Schedule) *entities.Schedule {
	schedule := &entities.Schedule{
		UUID:      scheduleModel.UUID,
		TenantID:  scheduleModel.TenantID,
		OwnerID:   scheduleModel.OwnerID,
		CreatedAt: scheduleModel.CreatedAt,
	}

	for _, job := range scheduleModel.Jobs {
		schedule.Jobs = append(schedule.Jobs, NewJobEntityFromModels(job))
	}

	return schedule
}

func NewScheduleModelFromEntities(schedule *entities.Schedule) *models.Schedule {
	scheduleModel := &models.Schedule{
		UUID:     schedule.UUID,
		TenantID: schedule.TenantID,
		OwnerID:  schedule.OwnerID,
	}

	for _, job := range schedule.Jobs {
		scheduleModel.Jobs = append(scheduleModel.Jobs, NewJobModelFromEntities(job, &scheduleModel.ID))
	}

	return scheduleModel
}
