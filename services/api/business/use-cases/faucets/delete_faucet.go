package faucets

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
)

const deleteFaucetComponent = "use-cases.delete-faucet"

// deleteFaucetUseCase is a use case to delete a faucet
type deleteFaucetUseCase struct {
	db store.DB
}

// NewDeleteFaucet creates a new DeleteFaucetUseCase
func NewDeleteFaucetUseCase(db store.DB) usecases.DeleteFaucetUseCase {
	return &deleteFaucetUseCase{
		db: db,
	}
}

// Execute deletes a faucet
func (uc *deleteFaucetUseCase) Execute(ctx context.Context, uuid string, tenants []string) error {
	logger := log.WithContext(ctx).WithField("faucet_uuid", uuid).WithField("tenants", tenants)
	logger.Debug("deleting faucet")

	faucetModel, err := uc.db.Faucet().FindOneByUUID(ctx, uuid, tenants)
	if err != nil {
		return errors.FromError(err).ExtendComponent(deleteFaucetComponent)
	}

	err = uc.db.Faucet().Delete(ctx, faucetModel, tenants)
	if err != nil {
		return errors.FromError(err).ExtendComponent(deleteFaucetComponent)
	}

	logger.Info("faucet deleted successfully")
	return nil
}
