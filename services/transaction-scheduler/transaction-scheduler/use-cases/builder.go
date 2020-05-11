package usecases

import (
	"github.com/Shopify/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/interfaces"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/jobs"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/schedules"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/transactions"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/validators"
)

type UseCases struct {
	SendTransaction transactions.SendTxUseCase
	CreateSchedule  schedules.CreateScheduleUseCase
	GetSchedule     schedules.GetScheduleUseCase
	CreateJob       jobs.CreateJobUseCase
	StartJob        jobs.StartJobUseCase
}

func NewUseCases(
	db interfaces.DB,
	chainRegistryClient client.ChainRegistryClient,
	producer sarama.SyncProducer,
	txCrafterPartition string,
) *UseCases {
	txValidator := validators.NewTransactionValidator(db, chainRegistryClient)

	// schedules
	createScheduleUseCase := schedules.NewCreateScheduleUseCase(txValidator, db)
	getScheduleUseCase := schedules.NewGetScheduleUseCase(db)

	// jobs
	createJobUseCase := jobs.NewCreateJobUseCase(db)
	startJobUseCase := jobs.NewStartJobUseCase(db, producer, txCrafterPartition)

	return &UseCases{
		SendTransaction: transactions.NewSendTxUseCase(txValidator, db, startJobUseCase),
		CreateSchedule:  createScheduleUseCase,
		GetSchedule:     getScheduleUseCase,
		CreateJob:       createJobUseCase,
		StartJob:        startJobUseCase,
	}
}
