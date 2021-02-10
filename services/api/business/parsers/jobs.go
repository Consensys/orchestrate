package parsers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models"
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
		Schedule: &models.Schedule{
			UUID:     job.ScheduleUUID,
			TenantID: job.TenantID,
		},
		Logs:      []*models.Log{},
		CreatedAt: job.CreatedAt,
		UpdatedAt: job.UpdatedAt,
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
	}

	if jobModel.Schedule != nil {
		job.ScheduleUUID = jobModel.Schedule.UUID
		job.TenantID = jobModel.Schedule.TenantID
	}

	if jobModel.Transaction != nil {
		job.Transaction = NewTransactionEntityFromModels(jobModel.Transaction)
	}

	lastLogID := -1
	for idx, logModel := range jobModel.Logs {
		job.Logs = append(job.Logs, NewLogEntityFromModels(logModel))
		// Ignore resending and warning statuses
		if logModel.Status == entities.StatusResending || logModel.Status == entities.StatusWarning {
			continue
		}
		// Ignore fail statuses if they come after a resending
		if logModel.Status == entities.StatusFailed && idx > 1 && jobModel.Logs[idx-1].Status == entities.StatusResending {
			continue
		}

		if logModel.ID > lastLogID {
			job.Status = logModel.Status
			lastLogID = logModel.ID
		}
	}

	return job
}
