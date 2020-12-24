package client

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/keymanager"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/keymanager/ethereum"

	healthz "github.com/heptiolabs/healthcheck"
	dto "github.com/prometheus/client_model/go"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
)

//go:generate mockgen -source=client.go -destination=mock/mock.go -package=mock

type OrchestrateClient interface {
	TransactionClient
	ScheduleClient
	JobClient
	MetricClient
	AccountClient
	FaucetClient
}

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

type AccountClient interface {
	CreateAccount(ctx context.Context, request *types.CreateAccountRequest) (*types.AccountResponse, error)
	SearchAccounts(ctx context.Context, filters *entities.AccountFilters) ([]*types.AccountResponse, error)
	GetAccount(ctx context.Context, address string) (*types.AccountResponse, error)
	ImportAccount(ctx context.Context, request *types.ImportAccountRequest) (*types.AccountResponse, error)
	UpdateAccount(ctx context.Context, address string, request *types.UpdateAccountRequest) (*types.AccountResponse, error)
	SignPayload(ctx context.Context, address string, request *types.SignPayloadRequest) (string, error)
	SignTypedData(ctx context.Context, address string, request *types.SignTypedDataRequest) (string, error)
	VerifySignature(ctx context.Context, request *keymanager.VerifyPayloadRequest) error
	VerifyTypedDataSignature(ctx context.Context, request *ethereum.VerifyTypedDataRequest) error
}

type FaucetClient interface {
	RegisterFaucet(ctx context.Context, request *types.RegisterFaucetRequest) (*types.FaucetResponse, error)
	UpdateFaucet(ctx context.Context, uuid string, request *types.UpdateFaucetRequest) (*types.FaucetResponse, error)
	GetFaucet(ctx context.Context, uuid string) (*types.FaucetResponse, error)
	SearchFaucets(ctx context.Context, filters *entities.FaucetFilters) ([]*types.FaucetResponse, error)
	DeleteFaucet(ctx context.Context, uuid string) error
}
