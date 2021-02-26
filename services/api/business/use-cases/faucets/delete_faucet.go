package faucets

import (
	"context"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
	usecases "github.com/ConsenSys/orchestrate/services/api/business/use-cases"
	"github.com/ConsenSys/orchestrate/services/api/store"
)

const deleteFaucetComponent = "use-cases.delete-faucet"

type deleteFaucetUseCase struct {
	db     store.DB
	logger *log.Logger
}

func NewDeleteFaucetUseCase(db store.DB) usecases.DeleteFaucetUseCase {
	return &deleteFaucetUseCase{
		db:     db,
		logger: log.NewLogger().SetComponent(deleteFaucetComponent),
	}
}

func (uc *deleteFaucetUseCase) Execute(ctx context.Context, uuid string, tenants []string) error {
	ctx = log.WithFields(ctx, log.Field("faucet", uuid))
	logger := uc.logger.WithContext(ctx)
	logger.Debug("deleting faucet")

	faucetModel, err := uc.db.Faucet().FindOneByUUID(ctx, uuid, tenants)
	if err != nil {
		return errors.FromError(err).ExtendComponent(deleteFaucetComponent)
	}

	err = uc.db.Faucet().Delete(ctx, faucetModel, tenants)
	if err != nil {
		return errors.FromError(err).ExtendComponent(deleteFaucetComponent)
	}

	logger.Info("faucet was deleted successfully")
	return nil
}
