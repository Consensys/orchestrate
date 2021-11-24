package faucets

import (
	"context"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
	"github.com/consensys/orchestrate/services/api/store"
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

func (uc *deleteFaucetUseCase) Execute(ctx context.Context, uuid string, userInfo *multitenancy.UserInfo) error {
	ctx = log.WithFields(ctx, log.Field("faucet", uuid))
	logger := uc.logger.WithContext(ctx)
	logger.Debug("deleting faucet")

	faucetModel, err := uc.db.Faucet().FindOneByUUID(ctx, uuid, userInfo.AllowedTenants)
	if err != nil {
		return errors.FromError(err).ExtendComponent(deleteFaucetComponent)
	}

	err = uc.db.Faucet().Delete(ctx, faucetModel, userInfo.AllowedTenants)
	if err != nil {
		return errors.FromError(err).ExtendComponent(deleteFaucetComponent)
	}

	logger.Info("faucet was deleted successfully")
	return nil
}
