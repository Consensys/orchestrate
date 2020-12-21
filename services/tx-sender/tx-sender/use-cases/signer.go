package usecases

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
)

//go:generate mockgen -source=signer.go -destination=mocks/signer.go -package=mocks

type SignETHTransactionUseCase interface {
	Execute(ctx context.Context, job *entities.Job) (signedRaw, txHash string, err error)
}

type SignEEATransactionUseCase interface {
	Execute(ctx context.Context, job *entities.Job) (raw, txHash string, err error)
}

type SignQuorumPrivateTransactionUseCase interface {
	Execute(ctx context.Context, job *entities.Job) (signedRaw, txHash string, err error)
}
