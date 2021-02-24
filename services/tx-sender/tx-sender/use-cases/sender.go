package usecases

import (
	"context"

	"github.com/ConsenSys/orchestrate/pkg/types/entities"
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
