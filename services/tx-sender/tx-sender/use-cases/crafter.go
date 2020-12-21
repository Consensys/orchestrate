package usecases

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
)

//go:generate mockgen -source=crafter.go -destination=mocks/crafter.go -package=mocks

type CraftTransactionUseCase interface {
	Execute(ctx context.Context, job *entities.Job) error
}
