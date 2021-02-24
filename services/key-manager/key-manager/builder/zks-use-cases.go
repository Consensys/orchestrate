package builder

import (
	usecases "github.com/ConsenSys/orchestrate/services/key-manager/key-manager/use-cases"
	zksnarks "github.com/ConsenSys/orchestrate/services/key-manager/key-manager/use-cases/zk-snarks"
	"github.com/ConsenSys/orchestrate/services/key-manager/store"
)

type zskUseCases struct {
	verifySignatureUC usecases.VerifyZKSSignatureUseCase
}

func NewZKSUseCases(_ store.Vault) usecases.ZKSUseCases {
	return &zskUseCases{
		verifySignatureUC: zksnarks.NewVerifySignatureUseCase(),
	}
}

func (u *zskUseCases) VerifySignature() usecases.VerifyZKSSignatureUseCase {
	return u.verifySignatureUC
}
