package jobs

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
	tsorm "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/orm"
)

//go:generate mockgen -source=create_job.go -destination=mocks/create_job.go -package=mocks

const createJobComponent = "use-cases.create-job"

type CreateJobUseCase interface {
	Execute(ctx context.Context, jobRequest *types.JobRequest, tenantID string) (*types.JobResponse, error)
}

// createJobUseCase is a use case to create a new transaction job
type createJobUseCase struct {
	db  store.DB
	orm tsorm.ORM
}

// NewCreateJobUseCase creates a new CreateJobUseCase
func NewCreateJobUseCase(db store.DB, orm tsorm.ORM) CreateJobUseCase {
	return &createJobUseCase{
		db:  db,
		orm: orm,
	}
}

// Execute validates and creates a new transaction job
func (uc *createJobUseCase) Execute(ctx context.Context, jobRequest *types.JobRequest, tenantID string) (*types.JobResponse, error) {
	log.WithContext(ctx).
		WithField("schedule_id", jobRequest.ScheduleUUID).
		WithField("tenant_id", tenantID).
		Debug("creating new job")

	err := utils.GetValidator().Struct(jobRequest)
	if err != nil {
		errMessage := "failed to validate create job request"
		log.WithError(err).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage).ExtendComponent(createJobComponent)
	}

	job, err := uc.buildJobFromRequest(ctx, jobRequest, tenantID)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createJobComponent)
	}

	if err := uc.orm.InsertOrUpdateJob(ctx, uc.db, job); err != nil {
		return nil, errors.FromError(err).ExtendComponent(createJobComponent)
	}

	log.WithContext(ctx).
		WithField("job_uuid", job.UUID).
		Info("job created successfully")

	return &types.JobResponse{
		UUID:        job.UUID,
		Transaction: jobRequest.Transaction,
		Status:      job.GetStatus(),
		CreatedAt:   job.CreatedAt,
	}, nil
}

func (uc *createJobUseCase) buildJobFromRequest(ctx context.Context, jobRequest *types.JobRequest, tenantID string) (*models.Job, error) {
	schedule, err := uc.db.Schedule().FindOneByUUID(ctx, jobRequest.ScheduleUUID, tenantID)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createJobComponent)
	}

	job := &models.Job{
		Type:       jobRequest.Type,
		Labels:     jobRequest.Labels,
		ScheduleID: &schedule.ID,
		Transaction: &models.Transaction{
			Hash:           jobRequest.Transaction.Hash,
			Sender:         jobRequest.Transaction.From,
			Recipient:      jobRequest.Transaction.To,
			Nonce:          jobRequest.Transaction.Nonce,
			Value:          jobRequest.Transaction.Value,
			GasPrice:       jobRequest.Transaction.GasPrice,
			GasLimit:       jobRequest.Transaction.GasLimit,
			Data:           jobRequest.Transaction.Data,
			PrivateFrom:    jobRequest.Transaction.PrivateFrom,
			PrivateFor:     jobRequest.Transaction.PrivateFor,
			PrivacyGroupID: jobRequest.Transaction.PrivacyGroupID,
			Raw:            jobRequest.Transaction.Raw,
		},
		Logs: []*models.Log{{
			Status:  types.JobStatusCreated,
			Message: "Job created",
		}},
	}

	return job, nil
}
