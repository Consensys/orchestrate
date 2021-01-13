package builder

import (
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases/transactions"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
)

type transactionUseCases struct {
	sendContractTransaction usecases.SendContractTxUseCase
	sendDeployTransaction   usecases.SendDeployTxUseCase
	sendTransaction         usecases.SendTxUseCase
	getTransaction          usecases.GetTxUseCase
	searchTransactions      usecases.SearchTransactionsUseCase
}

func newTransactionUseCases(
	db store.DB,
	searchChainsUC usecases.SearchChainsUseCase,
	getFaucetCandidateUC usecases.GetFaucetCandidateUseCase,
	schedulesUCs *scheduleUseCases,
	jobUCs *jobUseCases,
	getContractUC usecases.GetContractUseCase,
) *transactionUseCases {
	getTransactionUC := transactions.NewGetTxUseCase(db, schedulesUCs.GetSchedule())
	sendTxUC := transactions.NewSendTxUseCase(db, searchChainsUC, jobUCs.StartJob(), jobUCs.CreateJob(), schedulesUCs.CreateSchedule(), getTransactionUC, getFaucetCandidateUC)

	return &transactionUseCases{
		sendContractTransaction: transactions.NewSendContractTxUseCase(sendTxUC),
		sendDeployTransaction:   transactions.NewSendDeployTxUseCase(sendTxUC, getContractUC),
		sendTransaction:         sendTxUC,
		getTransaction:          getTransactionUC,
		searchTransactions:      transactions.NewSearchTransactionsUseCase(db, getTransactionUC),
	}
}

func (u *transactionUseCases) SendContractTransaction() usecases.SendContractTxUseCase {
	return u.sendContractTransaction
}

func (u *transactionUseCases) SendDeployTransaction() usecases.SendDeployTxUseCase {
	return u.sendDeployTransaction
}

func (u *transactionUseCases) SendTransaction() usecases.SendTxUseCase {
	return u.sendTransaction
}

func (u *transactionUseCases) GetTransaction() usecases.GetTxUseCase {
	return u.getTransaction
}

func (u *transactionUseCases) SearchTransactions() usecases.SearchTransactionsUseCase {
	return u.searchTransactions
}
