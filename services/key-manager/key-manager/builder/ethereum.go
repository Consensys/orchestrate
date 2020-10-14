package builder

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/key-manager/use-cases/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/store"
)

func NewEthereumUseCases(vault store.Vault) ethereum.UseCases {
	return &useCases{
		createAccount: ethereum.NewCreateAccountUseCase(vault),
	}
}

type useCases struct {
	createAccount ethereum.CreateAccountUseCase
}

func (ucs *useCases) CreateAccount() ethereum.CreateAccountUseCase {
	return ucs.createAccount
}
