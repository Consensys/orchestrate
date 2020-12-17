package builder

import (
	"github.com/Shopify/sarama"
	pkgsarama "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases/accounts"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases/jobs"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases/schedules"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases/transactions"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/validators"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/client"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/contract-registry/proto"
	keymanager "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/client"
)

type useCases struct {
	// Transactions
	sendContractTransaction usecases.SendContractTxUseCase
	sendDeployTransaction   usecases.SendDeployTxUseCase
	sendTransaction         usecases.SendTxUseCase
	getTransaction          usecases.GetTxUseCase
	searchTransactions      usecases.SearchTransactionsUseCase
	// Schedules
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
	// Accounts
	createAccountUC  usecases.CreateAccountUseCase
	getAccountUC     usecases.GetAccountUseCase
	searchAccountsUC usecases.SearchAccountsUseCase
	updateAccountUC  usecases.UpdateAccountUseCase
	fundAccountUC    usecases.FundAccountUseCase
}

func NewUseCases(
	db store.DB,
	appMetrics metrics.TransactionSchedulerMetrics,
	chainRegistryClient client.ChainRegistryClient,
	contractRegistryClient contractregistry.ContractRegistryClient,
	keyManagerClient keymanager.KeyManagerClient,
	producer sarama.SyncProducer,
	topicsCfg *pkgsarama.KafkaTopicConfig,
) usecases.UseCases {
	txValidator := validators.NewTransactionValidator(chainRegistryClient, contractRegistryClient)

	createScheduleUC := schedules.NewCreateScheduleUseCase(db)
	getScheduleUC := schedules.NewGetScheduleUseCase(db)
	createJobUC := jobs.NewCreateJobUseCase(db, txValidator)
	startJobUC := jobs.NewStartJobUseCase(db, producer, topicsCfg, appMetrics)
	resendJobUC := jobs.NewResendJobTxUseCase(db, producer, topicsCfg)
	updateChildrenUC := jobs.NewUpdateChildrenUseCase(db)
	startNextJobUC := jobs.NewStartNextJobUseCase(db, startJobUC)
	getTransactionUC := transactions.NewGetTxUseCase(db, getScheduleUC)

	sendTxUC := transactions.NewSendTxUseCase(txValidator, db, chainRegistryClient, startJobUC, createJobUC, createScheduleUC, getTransactionUC)
	searchAccountsUC := accounts.NewSearchAccountsUseCase(db)
	fundAccountUC := accounts.NewFundAccountUseCase(chainRegistryClient, sendTxUC)

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
		updateJob:   jobs.NewUpdateJobUseCase(db, updateChildrenUC, startNextJobUC, appMetrics),
		startJob:    startJobUC,
		resendJobTx: resendJobUC,
		// Accounts
		createAccountUC:  accounts.NewCreateAccountUseCase(db, searchAccountsUC, fundAccountUC, keyManagerClient),
		getAccountUC:     accounts.NewGetAccountUseCase(db),
		searchAccountsUC: searchAccountsUC,
		updateAccountUC:  accounts.NewUpdateAccountUseCase(db),
		fundAccountUC:    fundAccountUC,
	}
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

func (u *useCases) SearchJobs() usecases.SearchJobsUseCase {
	return u.searchJobs
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

func (u *useCases) GetAccount() usecases.GetAccountUseCase {
	return u.getAccountUC
}

func (u *useCases) SearchAccounts() usecases.SearchAccountsUseCase {
	return u.searchAccountsUC
}

func (u *useCases) CreateAccount() usecases.CreateAccountUseCase {
	return u.createAccountUC
}

func (u *useCases) UpdateAccount() usecases.UpdateAccountUseCase {
	return u.updateAccountUC
}

func (u *useCases) FundAccount() usecases.FundAccountUseCase {
	return u.fundAccountUC
}
