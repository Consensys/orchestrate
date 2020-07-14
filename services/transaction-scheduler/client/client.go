package client

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
)

//go:generate mockgen -source=client.go -destination=mock/mock.go -package=mock

type TransactionClient interface {
	SendContractTransaction(ctx context.Context, request *types.SendTransactionRequest) (*types.TransactionResponse, error)
	SendDeployTransaction(ctx context.Context, request *types.DeployContractRequest) (*types.TransactionResponse, error)
	SendRawTransaction(ctx context.Context, request *types.RawTransactionRequest) (*types.TransactionResponse, error)
	SendTransferTransaction(ctx context.Context, request *types.TransferRequest) (*types.TransactionResponse, error)
	GetTxRequest(ctx context.Context, txRequestUUID string) (*types.TransactionResponse, error)
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
	SearchJob(ctx context.Context, txHashes []string, chainUUID string) ([]*types.JobResponse, error)
}

type TransactionSchedulerClient interface {
	TransactionClient
	ScheduleClient
	JobClient
}
