package validators

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
)

type Validators struct {
	TransactionValidator TransactionValidator
}

func NewValidators(chainRegistryClient client.ChainRegistryClient) *Validators {
	return &Validators{
		TransactionValidator: NewTransaction(chainRegistryClient),
	}
}
