package client

import (
	"context"

	healthz "github.com/heptiolabs/healthcheck"
	dto "github.com/prometheus/client_model/go"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/txscheduler"
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
	ResendJobTx(ctx context.Context, jobUUID string) error
	SearchJob(ctx context.Context, filters *entities.JobFilters) ([]*types.JobResponse, error)
}

type MetricClient interface {
	Checker() healthz.Check
	Prometheus(context.Context) (map[string]*dto.MetricFamily, error)
}

type TransactionSchedulerClient interface {
	TransactionClient
	ScheduleClient
	JobClient
	MetricClient
}
