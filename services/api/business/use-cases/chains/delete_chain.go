package chains

import (
	"context"

	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/toolkit/database"
	"github.com/consensys/orchestrate/services/api/business/parsers"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
	"github.com/consensys/orchestrate/services/api/store"
)

const deleteChainComponent = "use-cases.delete-chain"

// deleteChainUseCase is a use case to delete a chain
type deleteChainUseCase struct {
	db         store.DB
	getChainUC usecases.GetChainUseCase
	logger     *log.Logger
}

// NewDeleteChainUseCase creates a new DeleteChainUseCase
func NewDeleteChainUseCase(db store.DB, getChainUC usecases.GetChainUseCase) usecases.DeleteChainUseCase {
	return &deleteChainUseCase{
		db:         db,
		getChainUC: getChainUC,
		logger:     log.NewLogger().SetComponent(deleteChainComponent),
	}
}

// Execute deletes a chain
func (uc *deleteChainUseCase) Execute(ctx context.Context, uuid string, userInfo *multitenancy.UserInfo) error {
	ctx = log.WithFields(ctx, log.Field("chain", uuid))
	logger := uc.logger.WithContext(ctx)
	logger.Debug("deleting chain")

	chain, err := uc.getChainUC.Execute(ctx, uuid, userInfo)
	if err != nil {
		return errors.FromError(err).ExtendComponent(deleteChainComponent)
	}

	chainModel := parsers.NewChainModelFromEntity(chain)
	chainModel.TenantID = userInfo.TenantID
	chainModel.OwnerID = userInfo.Username
	err = database.ExecuteInDBTx(uc.db, func(tx database.Tx) error {
		for _, privateTxManagerModel := range chainModel.PrivateTxManagers {
			der := tx.(store.Tx).PrivateTxManager().Delete(ctx, privateTxManagerModel)
			if der != nil {
				return der
			}
		}

		der := tx.(store.Tx).Chain().Delete(ctx, chainModel, userInfo.AllowedTenants)
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
