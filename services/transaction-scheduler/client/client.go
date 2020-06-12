package client

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/types"
)

//go:generate mockgen -source=client.go -destination=mock/mock.go -package=mock

type TransactionClient interface {
	SendContractTransaction(ctx context.Context, chainUUID string, request *types.SendTransactionRequest) (*types.TransactionResponse, error)
	SendDeployTransaction(ctx context.Context, chainUUID string, request *types.DeployContractRequest) (*types.TransactionResponse, error)
	SendRawTransaction(ctx context.Context, chainUUID string, request *types.RawTransactionRequest) (*types.TransactionResponse, error)
}

type ScheduleClient interface {
	GetSchedule(ctx context.Context, scheduleUUID string) (*types.ScheduleResponse, error)
	GetSchedules(ctx context.Context) ([]*types.ScheduleResponse, error)
	CreateSchedule(ctx context.Context, request *types.CreateScheduleRequest) (*types.ScheduleResponse, error)
}

type JobClient interface {
	GetJob(ctx context.Context, jobUUID string) (*types.JobResponse, error)
	GetJobs(ctx context.Context) ([]*types.JobResponse, error)
	CreateJob(ctx context.Context, request *types.CreateJobRequest) (*types.JobResponse, error)
	UpdateJob(ctx context.Context, jobUUID string, request *types.UpdateJobRequest) (*types.JobResponse, error)
	StartJob(ctx context.Context, jobUUID string) error
	SearchJob(ctx context.Context, txHashes []string) ([]*types.JobResponse, error)
}

type TransactionSchedulerClient interface {
	TransactionClient
	ScheduleClient
	JobClient
}
