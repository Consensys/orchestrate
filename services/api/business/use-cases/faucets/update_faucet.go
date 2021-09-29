package faucets

import (
	"context"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/services/api/business/parsers"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
	"github.com/consensys/orchestrate/services/api/store"
)

const updateFaucetComponent = "use-cases.update-faucet"

// updateFaucetUseCase is a use case to update a faucet
type updateFaucetUseCase struct {
	db     store.DB
	logger *log.Logger
}

// NewUpdateJobUseCase creates a new UpdateFaucetUseCase
func NewUpdateFaucetUseCase(db store.DB) usecases.UpdateFaucetUseCase {
	return &updateFaucetUseCase{
		db:     db,
		logger: log.NewLogger().SetComponent(updateFaucetComponent),
	}
}

// Execute updates a faucet
func (uc *updateFaucetUseCase) Execute(ctx context.Context, faucet *entities.Faucet, tenants []string) (*entities.Faucet, error) {
	ctx = log.WithFields(ctx, log.Field("faucet", faucet.UUID))
	logger := uc.logger.WithContext(ctx)
	logger.Debug("updating faucet")

	err := uc.db.Faucet().Update(ctx, parsers.NewFaucetModelFromEntity(faucet), tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(updateFaucetComponent)
	}

	faucetRetrieved, err := uc.db.Faucet().FindOneByUUID(ctx, faucet.UUID, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(updateFaucetComponent)
	}

	logger.Info("faucet updated successfully")
	return parsers.NewFaucetFromModel(faucetRetrieved), nil
}
