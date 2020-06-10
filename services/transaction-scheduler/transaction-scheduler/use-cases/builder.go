package usecases

import (
	"github.com/Shopify/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/jobs"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/schedules"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/transactions"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/validators"
)

type UseCases interface {
	transactions.UseCases
	schedules.UseCases
	jobs.UseCases
}

type useCases struct {
	// Transaction
	sendContractTransaction transactions.SendContractTxUseCase
	sendDeployTransaction   transactions.SendDeployTxUseCase
	// Schedule
	createSchedule schedules.CreateScheduleUseCase
	getSchedule    schedules.GetScheduleUseCase
	getSchedules   schedules.GetSchedulesUseCase
	// Jobs
	createJob  jobs.CreateJobUseCase
	getJob     jobs.GetJobUseCase
	startJob   jobs.StartJobUseCase
	updateJob  jobs.UpdateJobUseCase
	searchJobs jobs.SearchJobsUseCase
}

func NewUseCases(
	db store.DB,
	chainRegistryClient client.ChainRegistryClient,
	contractRegistryClient contractregistry.ContractRegistryClient,
	producer sarama.SyncProducer,
	txCrafterPartition string,
) UseCases {
	txValidator := validators.NewTransactionValidator(db, chainRegistryClient, contractRegistryClient)

	createScheduleUC := schedules.NewCreateScheduleUseCase(db)
	getScheduleUC := schedules.NewGetScheduleUseCase(db)
	createJobUC := jobs.NewCreateJobUseCase(db, txValidator)
	startJobUC := jobs.NewStartJobUseCase(db, producer, txCrafterPartition)
	sendTxUC := transactions.NewSendTxUseCase(txValidator, db, startJobUC, createJobUC, createScheduleUC, getScheduleUC)

	return &useCases{
		// Transaction
		sendContractTransaction: transactions.NewSendContractTxUseCase(txValidator, sendTxUC),
		sendDeployTransaction:   transactions.NewSendDeployTxUseCase(txValidator, sendTxUC),
		// Schedules
		createSchedule: createScheduleUC,
		getSchedule:    getScheduleUC,
		getSchedules:   schedules.NewGetSchedulesUseCase(db),
		// Jobs
		createJob:  createJobUC,
		getJob:     jobs.NewGetJobUseCase(db),
		searchJobs: jobs.NewSearchJobsUseCase(db),
		updateJob:  jobs.NewUpdateJobUseCase(db),
		startJob:   startJobUC,
	}
}

func (u *useCases) SendContractTransaction() transactions.SendContractTxUseCase {
	return u.sendContractTransaction
}

func (u *useCases) SendDeployTransaction() transactions.SendDeployTxUseCase {
	return u.sendDeployTransaction
}

func (u *useCases) CreateSchedule() schedules.CreateScheduleUseCase {
	return u.createSchedule
}

func (u *useCases) GetSchedule() schedules.GetScheduleUseCase {
	return u.getSchedule
}

func (u *useCases) GetSchedules() schedules.GetSchedulesUseCase {
	return u.getSchedules
}

func (u *useCases) CreateJob() jobs.CreateJobUseCase {
	return u.createJob
}

func (u *useCases) GetJob() jobs.GetJobUseCase {
	return u.getJob
}

func (u *useCases) StartJob() jobs.StartJobUseCase {
	return u.startJob
}

func (u *useCases) UpdateJob() jobs.UpdateJobUseCase {
	return u.updateJob
}

func (u *useCases) SearchJobs() jobs.SearchJobsUseCase {
	return u.searchJobs
}
