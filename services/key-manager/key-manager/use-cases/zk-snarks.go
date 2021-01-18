package usecases

import (
	"context"
)

//go:generate mockgen -source=zk-snarks.go -destination=mocks/zk-snarks.go -package=mocks

type ZKSUseCases interface {
	VerifySignature() VerifyZKSSignatureUseCase
}

type VerifyZKSSignatureUseCase interface {
	Execute(ctx context.Context, publicKey, signature, payload string) error
}
