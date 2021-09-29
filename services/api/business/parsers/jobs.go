package parsers

import (
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/services/api/store/models"
)

func NewJobModelFromEntities(job *entities.Job, scheduleID *int) *models.Job {
	jobModel := &models.Job{
		UUID:         job.UUID,
		ChainUUID:    job.ChainUUID,
		Type:         job.Type,
		NextJobUUID:  job.NextJobUUID,
		Labels:       job.Labels,
		InternalData: job.InternalData,
		ScheduleID:   scheduleID,
		Status:       job.Status,
		Schedule: &models.Schedule{
			UUID:     job.ScheduleUUID,
			TenantID: job.TenantID,
		},
		Logs:      []*models.Log{},
		CreatedAt: job.CreatedAt,
		UpdatedAt: job.UpdatedAt,
	}

	if job.InternalData != nil {
		jobModel.IsParent = job.InternalData.ParentJobUUID == ""
	}

	if job.Status == "" {
		job.Status = entities.StatusCreated
	}

	if scheduleID != nil {
		jobModel.Schedule.ID = *scheduleID
	}

	if job.Transaction != nil {
		jobModel.Transaction = NewTransactionModelFromEntities(job.Transaction)
	}

	for _, log := range job.Logs {
		jobModel.Logs = append(jobModel.Logs, NewLogModelFromEntity(log))
	}

	return jobModel
}

func NewJobEntityFromModels(jobModel *models.Job) *entities.Job {
	job := &entities.Job{
		UUID:         jobModel.UUID,
		ChainUUID:    jobModel.ChainUUID,
		NextJobUUID:  jobModel.NextJobUUID,
		Type:         jobModel.Type,
		Labels:       jobModel.Labels,
		InternalData: jobModel.InternalData,
		Logs:         []*entities.Log{},
		CreatedAt:    jobModel.CreatedAt,
		UpdatedAt:    jobModel.UpdatedAt,
		Status:       jobModel.Status,
	}

	if jobModel.Schedule != nil {
		job.ScheduleUUID = jobModel.Schedule.UUID
		job.TenantID = jobModel.Schedule.TenantID
	}

	if jobModel.Transaction != nil {
		job.Transaction = NewTransactionEntityFromModels(jobModel.Transaction)
	}

	for _, logModel := range jobModel.Logs {
		job.Logs = append(job.Logs, NewLogEntityFromModels(logModel))
	}

	return job
}
