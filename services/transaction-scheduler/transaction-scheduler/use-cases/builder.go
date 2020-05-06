package usecases

import (
	"github.com/Shopify/sarama"
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
	StartJob        jobs.StartJobUseCase
}

func NewUseCases(
	dataAgents *store.DataAgents,
	chainRegistryClient client.ChainRegistryClient,
	producer sarama.SyncProducer,
	txCrafterPartition string,
) *UseCases {
	txValidator := validators.NewTransactionValidator(dataAgents.TransactionRequest)

	createJobUseCase := jobs.NewCreateJobUseCase(dataAgents.JobAgent)
	startJobUseCase := jobs.NewStartJobUseCase(dataAgents.JobAgent, dataAgents.LogAgent, producer, txCrafterPartition)
	createScheduleUseCase := schedules.NewCreateScheduleUseCase(chainRegistryClient, dataAgents.ScheduleAgent)

	return &UseCases{
		SendTransaction: transactions.NewSendTxUseCase(
			txValidator,
			dataAgents.TransactionRequest,
			startJobUseCase,
		),
		CreateSchedule: createScheduleUseCase,
		CreateJob:      createJobUseCase,
		StartJob:       startJobUseCase,
	}
}
