package usecases

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
)

//go:generate mockgen -source=jobs.go -destination=mocks/jobs.go -package=mocks

/**
Jobs Use Cases
*/
type JobUseCases interface {
	CreateJob() CreateJobUseCase
	GetJob() GetJobUseCase
	StartJob() StartJobUseCase
	ResendJobTx() ResendJobTxUseCase
	UpdateJob() UpdateJobUseCase
	SearchJobs() SearchJobsUseCase
}

type CreateJobUseCase interface {
	Execute(ctx context.Context, job *entities.Job, tenants []string) (*entities.Job, error)
	WithDBTransaction(dbtx store.Tx) CreateJobUseCase
}

type GetJobUseCase interface {
	Execute(ctx context.Context, jobUUID string, tenants []string) (*entities.Job, error)
}

type SearchJobsUseCase interface {
	Execute(ctx context.Context, filters *entities.JobFilters, tenants []string) ([]*entities.Job, error)
}

type StartJobUseCase interface {
	Execute(ctx context.Context, jobUUID string, tenants []string) error
}

type StartNextJobUseCase interface {
	Execute(ctx context.Context, prevJobUUID string, tenants []string) error
}

type UpdateJobUseCase interface {
	Execute(ctx context.Context, jobEntity *entities.Job, nextStatus, logMessage string, tenants []string) (*entities.Job, error)
}

type UpdateChildrenUseCase interface {
	Execute(ctx context.Context, jobUUID, parentJobUUID, nextStatus string, tenants []string) error
	WithDBTransaction(dbtx store.Tx) UpdateChildrenUseCase
}

type ResendJobTxUseCase interface {
	Execute(ctx context.Context, jobUUID string, tenants []string) error
}
