package chains

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/parsers"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
)

const updateChainComponent = "use-cases.update-chain"

// updateChainUseCase is a use case to update a faucet
type updateChainUseCase struct {
	db         store.DB
	getChainUC usecases.GetChainUseCase
}

// NewUpdateChainUseCase creates a new UpdateChainUseCase
func NewUpdateChainUseCase(db store.DB, getChainUC usecases.GetChainUseCase) usecases.UpdateChainUseCase {
	return &updateChainUseCase{
		db:         db,
		getChainUC: getChainUC,
	}
}

// Execute updates a chain
func (uc *updateChainUseCase) Execute(ctx context.Context, chain *entities.Chain, tenants []string) (*entities.Chain, error) {
	logger := log.WithContext(ctx).WithField("chain_uuid", chain.UUID).WithField("tenants", tenants)
	logger.Debug("updating chain")

	chainRetrieved, err := uc.getChainUC.Execute(ctx, chain.UUID, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(updateChainComponent)
	}

	chainModel := parsers.NewChainModelFromEntity(chain)
	err = database.ExecuteInDBTx(uc.db, func(tx database.Tx) error {
		// If the chain has a private tx manager and we try to update it
		if chain.PrivateTxManager != nil && chainRetrieved.PrivateTxManager != nil {
			privateTxManager := chainModel.PrivateTxManagers[0]
			privateTxManager.UUID = chainRetrieved.PrivateTxManager.UUID
			der := tx.(store.Tx).PrivateTxManager().Update(ctx, privateTxManager)
			if der != nil {
				return der
			}
		}

		if chain.PrivateTxManager != nil && chainRetrieved.PrivateTxManager == nil {
			privateTxManager := chainModel.PrivateTxManagers[0]
			privateTxManager.ChainUUID = chainRetrieved.UUID
			der := tx.(store.Tx).PrivateTxManager().Insert(ctx, privateTxManager)
			if der != nil {
				return der
			}
		}

		der := tx.(store.Tx).Chain().Update(ctx, chainModel, tenants)
		if der != nil {
			return der
		}

		return nil
	})
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(updateChainComponent)
	}

	chainUpdated, err := uc.getChainUC.Execute(ctx, chain.UUID, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(updateChainComponent)
	}

	logger.Info("chain updated successfully")
	return chainUpdated, nil
}
