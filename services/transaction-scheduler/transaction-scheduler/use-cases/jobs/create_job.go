package jobs

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
)

// const createJobComponent = "use-cases.create-job"

//go:generate mockgen -source=create_job.go -destination=mocks/mock_create_job.go -package=mocks

type CreateJobUseCase interface {
	Execute(ctx context.Context, txJob *types.TransactionJob) (*types.TransactionJob, error)
}

// CreateJob is a use case to create a new transaction job
type CreateJob struct {
	txJobDataAgent store.TransactionJobAgent
}

// NewCreateJob creates a new CreateJob
func NewCreateJob(txJobDataAgent store.TransactionJobAgent) CreateJobUseCase {
	return &CreateJob{
		txJobDataAgent: txJobDataAgent,
	}
}

// Execute validates and creates a new transaction job in DB and posts a message on Kafka
func (usecase *CreateJob) Execute(ctx context.Context, txJob *types.TransactionJob) (*types.TransactionJob, error) {
	return nil, nil
}
