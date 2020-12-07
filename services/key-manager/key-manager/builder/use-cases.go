package builder

import (
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/key-manager/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/key-manager/use-cases/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/store"
)

type useCases struct {
	signTypedDataUC            usecases.SignTypedDataUseCase
	verifySignatureUC          usecases.VerifySignatureUseCase
	verifyTypedDataSignatureUC usecases.VerifyTypedDataSignatureUseCase
}

func NewUseCases(vaultClient store.Vault) usecases.UseCases {
	verifySignatureUC := ethereum.NewVerifySignatureUseCase()

	return &useCases{
		signTypedDataUC:            ethereum.NewSignTypedDataUseCase(vaultClient),
		verifySignatureUC:          verifySignatureUC,
		verifyTypedDataSignatureUC: ethereum.NewVerifyTypedDataSignatureUseCase(verifySignatureUC),
	}
}

func (u *useCases) SignTypedData() usecases.SignTypedDataUseCase {
	return u.signTypedDataUC
}

func (u *useCases) VerifySignature() usecases.VerifySignatureUseCase {
	return u.verifySignatureUC
}

func (u *useCases) VerifyTypedDataSignature() usecases.VerifyTypedDataSignatureUseCase {
	return u.verifyTypedDataSignatureUC
}
