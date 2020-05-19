package usecases

import (
	"github.com/Shopify/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/jobs"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/orm"
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
	sendTransaction transactions.SendTxUseCase
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
	producer sarama.SyncProducer,
	txCrafterPartition string,
) UseCases {
	txValidator := validators.NewTransactionValidator(db, chainRegistryClient)

	startJobUC := jobs.NewStartJobUseCase(db, producer, txCrafterPartition)

	return &useCases{
		// Transaction
		sendTransaction: transactions.NewSendTxUseCase(txValidator, db, orm.New(), startJobUC),
		// Schedules
		createSchedule: schedules.NewCreateScheduleUseCase(txValidator, db),
		getSchedule:    schedules.NewGetScheduleUseCase(db, orm.New()),
		getSchedules:   schedules.NewGetSchedulesUseCase(db, orm.New()),
		// Jobs
		createJob:  jobs.NewCreateJobUseCase(db, orm.New()),
		getJob:     jobs.NewGetJobUseCase(db),
		searchJobs: jobs.NewSearchJobsUseCase(db),
		updateJob:  jobs.NewUpdateJobUseCase(db, orm.New()),
		startJob:   startJobUC,
	}
}

func (u *useCases) SendTransaction() transactions.SendTxUseCase {
	return u.sendTransaction
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
