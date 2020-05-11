package jobs

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/interfaces"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
)

//go:generate mockgen -source=create_job.go -destination=mocks/create_job.go -package=mocks

const createJobComponent = "use-cases.create-job"

type CreateJobUseCase interface {
	Execute(ctx context.Context, jobRequest *types.JobRequest) (*types.JobResponse, error)
}

// createJobUseCase is a use case to create a new transaction job
type createJobUseCase struct {
	db interfaces.DB
}

// NewCreateJobUseCase creates a new CreateJobUseCase
func NewCreateJobUseCase(db interfaces.DB) CreateJobUseCase {
	return &createJobUseCase{db: db}
}

// Execute validates and creates a new transaction job
func (uc *createJobUseCase) Execute(ctx context.Context, jobRequest *types.JobRequest) (*types.JobResponse, error) {
	log.WithContext(ctx).WithField("schedule_id", jobRequest.ScheduleID).Debug("creating new job")

	err := utils.GetValidator().Struct(jobRequest)
	if err != nil {
		errMessage := "failed to validate job request"
		log.WithError(err).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage).ExtendComponent(createJobComponent)
	}

	job, err := uc.saveJob(ctx, jobRequest)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createJobComponent)
	}

	log.WithContext(ctx).WithField("job_uuid", job.UUID).Info("job created successfully")
	return &types.JobResponse{
		UUID:        job.UUID,
		Transaction: jobRequest.Transaction,
		Status:      job.GetStatus(),
		CreatedAt:   job.CreatedAt,
	}, nil
}

func (uc *createJobUseCase) saveJob(ctx context.Context, jobRequest *types.JobRequest) (*models.Job, error) {
	dbtx, err := uc.db.Begin()
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createJobComponent)
	}

	job := &models.Job{
		ScheduleID: jobRequest.ScheduleID,
		Type:       jobRequest.Type,
		Labels:     jobRequest.Labels,
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
	}
	err = dbtx.Job().Insert(ctx, job)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createJobComponent)
	}

	logModel := &models.Log{
		JobID:   job.ID,
		Status:  types.JobStatusCreated,
		Message: "Job created",
	}

	err = dbtx.Log().Insert(ctx, logModel)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createJobComponent)
	}

	err = dbtx.Commit()
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createJobComponent)
	}

	job.Logs = []*models.Log{logModel}
	return job, nil
}
