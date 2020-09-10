package builder

import (
	"github.com/Shopify/sarama"
	pkgsarama "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/jobs"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/schedules"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/transactions"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/validators"
)

func NewUseCases(
	db store.DB,
	chainRegistryClient client.ChainRegistryClient,
	contractRegistryClient contractregistry.ContractRegistryClient,
	producer sarama.SyncProducer,
	topicsCfg *pkgsarama.KafkaTopicConfig,
) usecases.UseCases {
	txValidator := validators.NewTransactionValidator(db, chainRegistryClient, contractRegistryClient)

	createScheduleUC := schedules.NewCreateScheduleUseCase(db)
	getScheduleUC := schedules.NewGetScheduleUseCase(db)
	createJobUC := jobs.NewCreateJobUseCase(db, txValidator)
	startJobUC := jobs.NewStartJobUseCase(db, producer, topicsCfg)
	resendJobUC := jobs.NewResendJobTxUseCase(db, producer, topicsCfg)
	updateChildrenUC := jobs.NewUpdateChildrenUseCase(db)
	startNextJobUC := jobs.NewStartNextJobUseCase(db, startJobUC)
	getTransactionUC := transactions.NewGetTxUseCase(db, getScheduleUC)

	sendTxUC := transactions.NewSendTxUseCase(txValidator, db, chainRegistryClient, startJobUC, createJobUC, createScheduleUC, getTransactionUC)

	return &useCases{
		// Transaction
		sendContractTransaction: transactions.NewSendContractTxUseCase(txValidator, sendTxUC),
		sendDeployTransaction:   transactions.NewSendDeployTxUseCase(txValidator, sendTxUC),
		sendTransaction:         sendTxUC,
		getTransaction:          getTransactionUC,
		searchTransactions:      transactions.NewSearchTransactionsUseCase(db, getTransactionUC),
		// Schedules
		createSchedule:  createScheduleUC,
		getSchedule:     getScheduleUC,
		searchSchedules: schedules.NewSearchSchedulesUseCase(db),
		// Jobs
		createJob:   createJobUC,
		getJob:      jobs.NewGetJobUseCase(db),
		searchJobs:  jobs.NewSearchJobsUseCase(db),
		updateJob:   jobs.NewUpdateJobUseCase(db, updateChildrenUC, startNextJobUC),
		startJob:    startJobUC,
		resendJobTx: resendJobUC,
	}
}

type useCases struct {
	// Transaction
	sendContractTransaction usecases.SendContractTxUseCase
	sendDeployTransaction   usecases.SendDeployTxUseCase
	sendTransaction         usecases.SendTxUseCase
	getTransaction          usecases.GetTxUseCase
	searchTransactions      usecases.SearchTransactionsUseCase
	// Schedule
	createSchedule  usecases.CreateScheduleUseCase
	getSchedule     usecases.GetScheduleUseCase
	searchSchedules usecases.SearchSchedulesUseCase
	// Jobs
	createJob   usecases.CreateJobUseCase
	getJob      usecases.GetJobUseCase
	startJob    usecases.StartJobUseCase
	resendJobTx usecases.ResendJobTxUseCase
	updateJob   usecases.UpdateJobUseCase
	searchJobs  usecases.SearchJobsUseCase
}

func (u *useCases) SendContractTransaction() usecases.SendContractTxUseCase {
	return u.sendContractTransaction
}

func (u *useCases) SendDeployTransaction() usecases.SendDeployTxUseCase {
	return u.sendDeployTransaction
}

func (u *useCases) SendTransaction() usecases.SendTxUseCase {
	return u.sendTransaction
}

func (u *useCases) GetTransaction() usecases.GetTxUseCase {
	return u.getTransaction
}

func (u *useCases) SearchTransactions() usecases.SearchTransactionsUseCase {
	return u.searchTransactions
}

func (u *useCases) CreateSchedule() usecases.CreateScheduleUseCase {
	return u.createSchedule
}

func (u *useCases) GetSchedule() usecases.GetScheduleUseCase {
	return u.getSchedule
}

func (u *useCases) SearchSchedules() usecases.SearchSchedulesUseCase {
	return u.searchSchedules
}

func (u *useCases) CreateJob() usecases.CreateJobUseCase {
	return u.createJob
}

func (u *useCases) GetJob() usecases.GetJobUseCase {
	return u.getJob
}

func (u *useCases) StartJob() usecases.StartJobUseCase {
	return u.startJob
}

func (u *useCases) ResendJobTx() usecases.ResendJobTxUseCase {
	return u.resendJobTx
}

func (u *useCases) UpdateJob() usecases.UpdateJobUseCase {
	return u.updateJob
}

func (u *useCases) SearchJobs() usecases.SearchJobsUseCase {
	return u.searchJobs
}
