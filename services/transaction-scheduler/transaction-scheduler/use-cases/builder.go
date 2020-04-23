package usecases

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/transactions"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/validators"
)

type UseCases struct {
	SendTransaction transactions.SendTxUseCase
}

func NewUseCases(dataAgents *store.DataAgents, val *validators.Validators) *UseCases {
	return &UseCases{
		SendTransaction: transactions.NewSendTx(dataAgents.TransactionRequest, val.TransactionValidator),
	}
}
