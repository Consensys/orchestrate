package client

import (
	"context"

	qkmtypes "github.com/consensys/quorum-key-manager/src/stores/api/types"

	types "github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/pkg/types/entities"
	ethcommon "github.com/ethereum/go-ethereum/common"
	healthz "github.com/heptiolabs/healthcheck"
	dto "github.com/prometheus/client_model/go"
)

//go:generate mockgen -source=client.go -destination=mock/mock.go -package=mock

type OrchestrateClient interface {
	TransactionClient
	ScheduleClient
	JobClient
	MetricClient
	AccountClient
	FaucetClient
	ChainClient
	ContractClient
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
	GetAccount(ctx context.Context, address ethcommon.Address) (*types.AccountResponse, error)
	ImportAccount(ctx context.Context, request *types.ImportAccountRequest) (*types.AccountResponse, error)
	UpdateAccount(ctx context.Context, address ethcommon.Address, request *types.UpdateAccountRequest) (*types.AccountResponse, error)
	SignMessage(ctx context.Context, address ethcommon.Address, request *qkmtypes.SignMessageRequest) (string, error)
	SignTypedData(ctx context.Context, address ethcommon.Address, request *qkmtypes.SignTypedDataRequest) (string, error)
	VerifyMessageSignature(ctx context.Context, request *qkmtypes.VerifyRequest) error
	VerifyTypedDataSignature(ctx context.Context, request *qkmtypes.VerifyTypedDataRequest) error
}

type FaucetClient interface {
	RegisterFaucet(ctx context.Context, request *types.RegisterFaucetRequest) (*types.FaucetResponse, error)
	UpdateFaucet(ctx context.Context, uuid string, request *types.UpdateFaucetRequest) (*types.FaucetResponse, error)
	GetFaucet(ctx context.Context, uuid string) (*types.FaucetResponse, error)
	SearchFaucets(ctx context.Context, filters *entities.FaucetFilters) ([]*types.FaucetResponse, error)
	DeleteFaucet(ctx context.Context, uuid string) error
}

type ChainClient interface {
	RegisterChain(ctx context.Context, request *types.RegisterChainRequest) (*types.ChainResponse, error)
	UpdateChain(ctx context.Context, uuid string, request *types.UpdateChainRequest) (*types.ChainResponse, error)
	GetChain(ctx context.Context, uuid string) (*types.ChainResponse, error)
	SearchChains(ctx context.Context, filters *entities.ChainFilters) ([]*types.ChainResponse, error)
	DeleteChain(ctx context.Context, uuid string) error
}

type ContractClient interface {
	RegisterContract(ctx context.Context, req *types.RegisterContractRequest) (*types.ContractResponse, error)
	DeregisterContract(ctx context.Context, name, tag string) error
	GetContract(ctx context.Context, name, tag string) (*types.ContractResponse, error)
	GetContractsCatalog(ctx context.Context) ([]string, error)
	GetContractTags(ctx context.Context, name string) ([]string, error)
	SetContractAddressCodeHash(ctx context.Context, address, chainID string, req *types.SetContractCodeHashRequest) error
	GetContractEvents(ctx context.Context, address, chainID string, req *types.GetContractEventsRequest) (*types.GetContractEventsBySignHashResponse, error)
	GetContractMethodSignatures(ctx context.Context, name, tag, method string) ([]string, error)
}
