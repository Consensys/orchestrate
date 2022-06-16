package builder

import (
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
	"github.com/consensys/orchestrate/services/api/business/use-cases/transactions"
	"github.com/consensys/orchestrate/services/api/store"
)

type transactionUseCases struct {
	sendContractTransaction usecases.SendContractTxUseCase
	sendDeployTransaction   usecases.SendDeployTxUseCase
	sendTransaction         usecases.SendTxUseCase
	getTransaction          usecases.GetTxUseCase
	searchTransactions      usecases.SearchTransactionsUseCase
	speedUp                 usecases.SpeedUpTxUseCase
	callOff                 usecases.CallOffTxUseCase
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
	sendTxUC := transactions.NewSendTxUseCase(db, searchChainsUC, jobUCs.StartJob(), jobUCs.CreateJob(), getTransactionUC, getFaucetCandidateUC)

	return &transactionUseCases{
		sendContractTransaction: transactions.NewSendContractTxUseCase(sendTxUC, getContractUC),
		sendDeployTransaction:   transactions.NewSendDeployTxUseCase(sendTxUC, getContractUC),
		sendTransaction:         sendTxUC,
		getTransaction:          getTransactionUC,
		searchTransactions:      transactions.NewSearchTransactionsUseCase(db, getTransactionUC),
		speedUp:                 transactions.NewSpeedUpTxUseCase(getTransactionUC, jobUCs.RetryTx()),
		callOff:                 transactions.NewCallOffTxUseCase(getTransactionUC, jobUCs.RetryTx()),
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

func (u *transactionUseCases) SpeedUp() usecases.SpeedUpTxUseCase {
	return u.speedUp
}

func (u *transactionUseCases) CallOff() usecases.CallOffTxUseCase {
	return u.callOff
}
