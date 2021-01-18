package builder

import (
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/key-manager/use-cases"
	zksnarks "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/key-manager/use-cases/zk-snarks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/store"
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
