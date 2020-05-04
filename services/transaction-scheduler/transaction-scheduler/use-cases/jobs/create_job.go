package jobs

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
)

//go:generate mockgen -source=create_job.go -destination=mocks/create_job.go -package=mocks

const createJobComponent = "use-cases.create-job"

type CreateJobUseCase interface {
	Execute(ctx context.Context, jobRequest *types.JobRequest) (*types.JobResponse, error)
}

// createJob is a use case to create a new transaction job
type createJob struct {
	jobDataAgent store.JobAgent
}

// NewCreateJob creates a new CreateJobUseCase
func NewCreateJob(jobDataAgent store.JobAgent) CreateJobUseCase {
	return &createJob{jobDataAgent: jobDataAgent}
}

// Execute validates and creates a new transaction job
func (usecase *createJob) Execute(ctx context.Context, jobRequest *types.JobRequest) (*types.JobResponse, error) {
	log.WithContext(ctx).WithField("schedule_id", jobRequest.ScheduleID).Debug("creating new job")

	// TODO: Add validation when use case becomes externally available through API

	job := &models.Job{
		ScheduleID: jobRequest.ScheduleID,
		Type:       jobRequest.Type,
		Labels:     jobRequest.Labels,
		Transaction: &models.Transaction{
			Hash:           *jobRequest.Transaction.Hash,
			Sender:         *jobRequest.Transaction.From,
			Recipient:      *jobRequest.Transaction.To,
			Nonce:          *jobRequest.Transaction.Nonce,
			Value:          *jobRequest.Transaction.Value,
			GasPrice:       *jobRequest.Transaction.GasPrice,
			GasLimit:       *jobRequest.Transaction.GasLimit,
			Data:           *jobRequest.Transaction.Data,
			PrivateFrom:    *jobRequest.Transaction.PrivateFrom,
			PrivateFor:     *jobRequest.Transaction.PrivateFor,
			PrivacyGroupID: *jobRequest.Transaction.PrivacyGroupID,
			Raw:            *jobRequest.Transaction.Raw,
		},
		Logs: []*models.Log{{
			Status: types.LogStatusCreated,
		}},
	}
	err := usecase.jobDataAgent.Insert(ctx, job)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createJobComponent)
	}

	log.WithContext(ctx).WithField("job_uuid", job.UUID).Info("job created successfully")
	return &types.JobResponse{
		UUID:        job.UUID,
		Transaction: jobRequest.Transaction,
		CreatedAt:   job.CreatedAt,
	}, nil
}
