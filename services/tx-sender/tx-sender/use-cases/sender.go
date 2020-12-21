package usecases

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
)

//go:generate mockgen -source=sender.go -destination=mocks/sender.go -package=mocks

type SendETHRawTxUseCase interface {
	Execute(ctx context.Context, job *entities.Job) error
}

type SendETHTxUseCase interface {
	Execute(ctx context.Context, job *entities.Job) error
}

type SendEEAPrivateTxUseCase interface {
	Execute(ctx context.Context, job *entities.Job) error
}

type SendTesseraPrivateTxUseCase interface {
	Execute(ctx context.Context, job *entities.Job) error
}

type SendTesseraMarkingTxUseCase interface {
	Execute(ctx context.Context, job *entities.Job) error
}
