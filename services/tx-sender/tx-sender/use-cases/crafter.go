package usecases

import (
	"context"

	"github.com/consensys/orchestrate/pkg/types/entities"
)

//go:generate mockgen -source=crafter.go -destination=mocks/crafter.go -package=mocks

type CraftTransactionUseCase interface {
	Execute(ctx context.Context, job *entities.Job) error
}
