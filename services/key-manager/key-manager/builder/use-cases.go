package builder

import (
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/key-manager/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/key-manager/use-cases/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/store"
)

type useCases struct {
	signTypedData usecases.SignTypedDataUseCase
}

func NewUseCases(vaultClient store.Vault) usecases.UseCases {
	return &useCases{
		signTypedData: ethereum.NewSignTypedDataUseCase(vaultClient),
	}
}

func (u *useCases) SignTypedData() usecases.SignTypedDataUseCase {
	return u.signTypedData
}
