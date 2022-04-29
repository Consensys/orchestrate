package accounts

import (
	"context"

	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
	"github.com/consensys/quorum-key-manager/pkg/client"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/services/api/store"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const deleteAccountComponent = "use-cases.delete-account"

type deleteAccountUseCase struct {
	db               store.DB
	keyManagerClient client.EthClient
	logger           *log.Logger
}

func NewDeleteAccountUseCase(db store.DB, keyManagerClient client.EthClient) usecases.DeleteAccountUseCase {
	return &deleteAccountUseCase{
		db:               db,
		keyManagerClient: keyManagerClient,
		logger:           log.NewLogger().SetComponent(deleteAccountComponent),
	}
}

func (uc *deleteAccountUseCase) Execute(ctx context.Context, address ethcommon.Address, userInfo *multitenancy.UserInfo) error {
	ctx = log.WithFields(ctx, log.Field("address", address))
	logger := uc.logger.WithContext(ctx)

	model, err := uc.db.Account().FindOneByAddress(ctx, address.Hex(), userInfo.AllowedTenants, userInfo.Username)
	if err != nil {
		return errors.FromError(err).ExtendComponent(deleteAccountComponent)
	}

	err = uc.db.Account().Delete(ctx, address.Hex(), userInfo.AllowedTenants, userInfo.Username)
	if err != nil {
		return errors.FromError(err).ExtendComponent(deleteAccountComponent)
	}

	// First, we soft delete the account
	err = uc.keyManagerClient.DeleteEthAccount(ctx, model.StoreID, model.Address)
	if err != nil {
		errMsg := "failed to remove quorum key manager account"
		uc.logger.WithError(err).Error(errMsg)
		return errors.DependencyFailureError(errMsg)
	}

	// Then, we destroy it
	err = uc.keyManagerClient.DestroyEthAccount(ctx, model.StoreID, model.Address)
	if err != nil {
		errMsg := "failed to destroy quorum key manager account"
		uc.logger.WithError(err).Error(errMsg)
		return errors.DependencyFailureError(errMsg)
	}

	logger.Info("account deleted successfully")
	return nil
}
