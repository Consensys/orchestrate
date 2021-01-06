package builder

import (
	"github.com/Shopify/sarama"
	pkgsarama "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases/accounts"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases/contracts"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases/faucets"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases/jobs"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases/schedules"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases/transactions"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/validators"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/client"
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

	// Faucets
	registerFaucetUC usecases.RegisterFaucetUseCase
	updateFaucetUC   usecases.UpdateFaucetUseCase
	getFaucetUC      usecases.GetFaucetUseCase
	searchFaucetUC   usecases.SearchFaucetsUseCase
	deleteFaucetUC   usecases.DeleteFaucetUseCase

	// Contracts
	GetContractsCatalogUC        usecases.GetContractsCatalogUseCase
	getContractEvents           usecases.GetContractEventsUseCase
	getContractMethodSignatures usecases.GetContractMethodSignaturesUseCase
	getContractMethods          usecases.GetContractMethodsUseCase
	getContractTags             usecases.GetContractTagsUseCase
	setContractCodeHash         usecases.SetContractCodeHashUseCase
	registerContractUC          usecases.RegisterContractUseCase
	getContractUC               usecases.GetContractUseCase
}

func NewUseCases(
	db store.DB,
	appMetrics metrics.TransactionSchedulerMetrics,
	chainRegistryClient client.ChainRegistryClient,
	keyManagerClient keymanager.KeyManagerClient,
	chainStateReader ethclient.ChainStateReader,
	producer sarama.SyncProducer,
	topicsCfg *pkgsarama.KafkaTopicConfig,
) usecases.UseCases {
	getContractUC := contracts.NewGetContractUseCase(db.Artifact())

	txValidator := validators.NewTransactionValidator(chainRegistryClient)

	createScheduleUC := schedules.NewCreateScheduleUseCase(db)
	getScheduleUC := schedules.NewGetScheduleUseCase(db)
	createJobUC := jobs.NewCreateJobUseCase(db, txValidator)
	startJobUC := jobs.NewStartJobUseCase(db, producer, topicsCfg, appMetrics)
	resendJobUC := jobs.NewResendJobTxUseCase(db, producer, topicsCfg)
	updateChildrenUC := jobs.NewUpdateChildrenUseCase(db)
	startNextJobUC := jobs.NewStartNextJobUseCase(db, startJobUC)
	getTransactionUC := transactions.NewGetTxUseCase(db, getScheduleUC)
	searchFaucetsUC := faucets.NewSearchFaucets(db)
	getFaucetCandidate := faucets.NewGetFaucetCandidateUseCase(chainRegistryClient, searchFaucetsUC, chainStateReader)

	sendTxUC := transactions.NewSendTxUseCase(txValidator, db, chainRegistryClient, startJobUC, createJobUC, createScheduleUC, getTransactionUC, getFaucetCandidate)
	searchAccountsUC := accounts.NewSearchAccountsUseCase(db)
	fundAccountUC := accounts.NewFundAccountUseCase(chainRegistryClient, sendTxUC, getFaucetCandidate)

	return &useCases{
		// Transaction
		sendContractTransaction: transactions.NewSendContractTxUseCase(sendTxUC),
		sendDeployTransaction:   transactions.NewSendDeployTxUseCase(sendTxUC, getContractUC),
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

		// Faucets
		registerFaucetUC: faucets.NewRegisterFaucetUseCase(db, searchFaucetsUC),
		updateFaucetUC:   faucets.NewUpdateFaucetUseCase(db),
		getFaucetUC:      faucets.NewGetFaucetUseCase(db),
		searchFaucetUC:   searchFaucetsUC,
		deleteFaucetUC:   faucets.NewDeleteFaucetUseCase(db),

		// Contracts
		registerContractUC:          contracts.NewRegisterContractUseCase(db),
		getContractUC:               getContractUC,
		GetContractsCatalogUC:        contracts.NewGetCatalogUseCase(db.Repository()),
		getContractEvents:           contracts.NewGetEventsUseCase(db.Event()),
		getContractMethodSignatures: contracts.NewGetMethodSignaturesUseCase(getContractUC),
		getContractMethods:          contracts.NewGetMethodsUseCase(db.Method()),
		getContractTags:             contracts.NewGetTagsUseCase(db.Tag()),
		setContractCodeHash:         contracts.NewSetCodeHashUseCase(db.CodeHash()),
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

func (u *useCases) RegisterFaucet() usecases.RegisterFaucetUseCase {
	return u.registerFaucetUC
}

func (u *useCases) UpdateFaucet() usecases.UpdateFaucetUseCase {
	return u.updateFaucetUC
}

func (u *useCases) GetFaucet() usecases.GetFaucetUseCase {
	return u.getFaucetUC
}

func (u *useCases) SearchFaucets() usecases.SearchFaucetsUseCase {
	return u.searchFaucetUC
}

func (u *useCases) DeleteFaucet() usecases.DeleteFaucetUseCase {
	return u.deleteFaucetUC
}

func (u *useCases) GetContract() usecases.GetContractUseCase {
	return u.getContractUC
}

func (u *useCases) RegisterContract() usecases.RegisterContractUseCase {
	return u.registerContractUC
}

func (u *useCases) GetContractsCatalog() usecases.GetContractsCatalogUseCase {
	return u.GetContractsCatalogUC
}

func (u *useCases) GetContractEvents() usecases.GetContractEventsUseCase {
	return u.getContractEvents
}

func (u *useCases) GetContractMethodSignatures() usecases.GetContractMethodSignaturesUseCase {
	return u.getContractMethodSignatures
}

func (u *useCases) GetContractMethods() usecases.GetContractMethodsUseCase {
	return u.getContractMethods
}

func (u *useCases) GetContractTags() usecases.GetContractTagsUseCase {
	return u.getContractTags
}

func (u *useCases) SetContractCodeHash() usecases.SetContractCodeHashUseCase {
	return u.setContractCodeHash
}
