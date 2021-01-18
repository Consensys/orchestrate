package builder

import (
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/key-manager/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/key-manager/use-cases/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/store"
)

type ethUseCases struct {
	signTypedDataUC            usecases.SignTypedDataUseCase
	verifySignatureUC          usecases.VerifyETHSignatureUseCase
	verifyTypedDataSignatureUC usecases.VerifyTypedDataSignatureUseCase
}

func NewETHUseCases(vaultClient store.Vault) usecases.ETHUseCases {
	verifySignatureUC := ethereum.NewVerifySignatureUseCase()

	return &ethUseCases{
		signTypedDataUC:            ethereum.NewSignTypedDataUseCase(vaultClient),
		verifySignatureUC:          verifySignatureUC,
		verifyTypedDataSignatureUC: ethereum.NewVerifyTypedDataSignatureUseCase(verifySignatureUC),
	}
}

func (u *ethUseCases) SignTypedData() usecases.SignTypedDataUseCase {
	return u.signTypedDataUC
}

func (u *ethUseCases) VerifySignature() usecases.VerifyETHSignatureUseCase {
	return u.verifySignatureUC
}

func (u *ethUseCases) VerifyTypedDataSignature() usecases.VerifyTypedDataSignatureUseCase {
	return u.verifyTypedDataSignatureUC
}
