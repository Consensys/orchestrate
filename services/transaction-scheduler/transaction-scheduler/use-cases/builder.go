package usecases

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/jobs"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/schedules"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/transactions"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/validators"
)

type UseCases struct {
	SendTransaction transactions.SendTxUseCase
	CreateSchedule  schedules.CreateScheduleUseCase
	CreateJob       jobs.CreateJobUseCase
}

func NewUseCases(dataAgents *store.DataAgents, chainRegistryClient client.ChainRegistryClient) *UseCases {
	vals := validators.NewValidators(dataAgents.TransactionRequest)

	createJobUseCase := jobs.NewCreateJob(dataAgents.JobAgent)
	createScheduleUseCase := schedules.NewCreateSchedule(chainRegistryClient, dataAgents.ScheduleAgent)

	return &UseCases{
		SendTransaction: transactions.NewSendTx(dataAgents.TransactionRequest, vals.TransactionValidator, createScheduleUseCase, createJobUseCase),
		CreateSchedule:  createScheduleUseCase,
		CreateJob:       createJobUseCase,
	}
}
