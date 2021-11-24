package usecases

import (
	"context"

	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/services/api/store"
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
	Execute(ctx context.Context, job *entities.Job, userInfo *multitenancy.UserInfo) (*entities.Job, error)
	WithDBTransaction(dbtx store.Tx) CreateJobUseCase
}

type GetJobUseCase interface {
	Execute(ctx context.Context, jobUUID string, userInfo *multitenancy.UserInfo) (*entities.Job, error)
}

type SearchJobsUseCase interface {
	Execute(ctx context.Context, filters *entities.JobFilters, userInfo *multitenancy.UserInfo) ([]*entities.Job, error)
}

type StartJobUseCase interface {
	Execute(ctx context.Context, jobUUID string, userInfo *multitenancy.UserInfo) error
}

type StartNextJobUseCase interface {
	Execute(ctx context.Context, prevJobUUID string, userInfo *multitenancy.UserInfo) error
}

type UpdateJobUseCase interface {
	Execute(ctx context.Context, jobEntity *entities.Job, nextStatus entities.JobStatus, logMessage string, userInfo *multitenancy.UserInfo) (*entities.Job, error)
}

type UpdateChildrenUseCase interface {
	Execute(ctx context.Context, jobUUID, parentJobUUID string, nextStatus entities.JobStatus, userInfo *multitenancy.UserInfo) error
	WithDBTransaction(dbtx store.Tx) UpdateChildrenUseCase
}

type ResendJobTxUseCase interface {
	Execute(ctx context.Context, jobUUID string, userInfo *multitenancy.UserInfo) error
}
