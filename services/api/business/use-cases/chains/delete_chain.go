package chains

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/parsers"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
)

const deleteChainComponent = "use-cases.delete-chain"

// deleteChainUseCase is a use case to delete a chain
type deleteChainUseCase struct {
	db         store.DB
	getChainUC usecases.GetChainUseCase
}

// NewDeleteChainUseCase creates a new DeleteChainUseCase
func NewDeleteChainUseCase(db store.DB, getChainUC usecases.GetChainUseCase) usecases.DeleteChainUseCase {
	return &deleteChainUseCase{
		db:         db,
		getChainUC: getChainUC,
	}
}

// Execute deletes a chain
func (uc *deleteChainUseCase) Execute(ctx context.Context, uuid string, tenants []string) error {
	logger := log.WithContext(ctx).WithField("chain_uuid", uuid).WithField("tenants", tenants)
	logger.Debug("deleting chain")

	chain, err := uc.getChainUC.Execute(ctx, uuid, tenants)
	if err != nil {
		return errors.FromError(err).ExtendComponent(deleteChainComponent)
	}

	chainModel := parsers.NewChainModelFromEntity(chain)
	err = database.ExecuteInDBTx(uc.db, func(tx database.Tx) error {
		for _, privateTxManagerModel := range chainModel.PrivateTxManagers {
			der := tx.(store.Tx).PrivateTxManager().Delete(ctx, privateTxManagerModel)
			if der != nil {
				return der
			}
		}

		der := tx.(store.Tx).Chain().Delete(ctx, chainModel, tenants)
		if der != nil {
			return der
		}

		return nil
	})
	if err != nil {
		return errors.FromError(err).ExtendComponent(updateChainComponent)
	}

	logger.Info("chain deleted successfully")
	return nil
}
