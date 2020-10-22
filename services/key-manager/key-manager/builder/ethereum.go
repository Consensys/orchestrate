package builder

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/key-manager/use-cases/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/store"
)

type useCases struct {
	createAccount ethereum.CreateAccountUseCase
	sign          ethereum.SignUseCase
	signTx        ethereum.SignTransactionUseCase
	signTesseraTx ethereum.SignTesseraTransactionUseCase
}

func NewEthereumUseCases(vault store.Vault) ethereum.UseCases {
	return &useCases{
		createAccount: ethereum.NewCreateAccountUseCase(vault),
		sign:          ethereum.NewSignUseCase(vault),
		signTx:        ethereum.NewSignTransactionUseCase(vault),
		signTesseraTx: ethereum.NewSignTesseraTransactionUseCase(vault),
	}
}

func (ucs *useCases) CreateAccount() ethereum.CreateAccountUseCase {
	return ucs.createAccount
}

func (ucs *useCases) SignPayload() ethereum.SignUseCase {
	return ucs.sign
}

func (ucs *useCases) SignTransaction() ethereum.SignTransactionUseCase {
	return ucs.signTx
}

func (ucs *useCases) SignTesseraTransaction() ethereum.SignTesseraTransactionUseCase {
	return ucs.signTesseraTx
}
