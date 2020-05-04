package validators

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
)

type Validators struct {
	TransactionValidator TransactionValidator
}

func NewValidators(txRequestAgent store.TransactionRequestAgent) *Validators {
	return &Validators{
		TransactionValidator: NewTransactionValidator(txRequestAgent),
	}
}
