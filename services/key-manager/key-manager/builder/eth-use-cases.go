package builder

import (
	usecases "github.com/ConsenSys/orchestrate/services/key-manager/key-manager/use-cases"
	"github.com/ConsenSys/orchestrate/services/key-manager/key-manager/use-cases/ethereum"
	"github.com/ConsenSys/orchestrate/services/key-manager/store"
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
