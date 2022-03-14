package faucets

import (
	"context"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/services/api/business/parsers"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
	"github.com/consensys/orchestrate/services/api/store"
	ethcommon "github.com/ethereum/go-ethereum/common"
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
func (uc *updateFaucetUseCase) Execute(ctx context.Context, faucet *entities.Faucet, userInfo *multitenancy.UserInfo) (*entities.Faucet, error) {
	ctx = log.WithFields(ctx, log.Field("faucet", faucet.UUID))
	logger := uc.logger.WithContext(ctx)
	logger.Debug("updating faucet")

	if faucet.ChainRule != "" {
		_, err := uc.db.Chain().FindOneByUUID(ctx, faucet.ChainRule, userInfo.AllowedTenants, userInfo.Username)
		if errors.IsNotFoundError(err) {
			return nil, errors.InvalidParameterError("cannot find new linked chain")
		} else if err != nil {
			return nil, err
		}
	}

	if faucet.CreditorAccount.String() != new(ethcommon.Address).String() {
		_, err := uc.db.Account().FindOneByAddress(ctx, faucet.CreditorAccount.String(), userInfo.AllowedTenants, userInfo.Username)
		if errors.IsNotFoundError(err) {
			return nil, errors.InvalidParameterError("cannot find updated creditor account")
		} else if err != nil {
			return nil, err
		}
	}

	faucetModel := parsers.NewFaucetModelFromEntity(faucet)
	faucetModel.TenantID = userInfo.TenantID
	err := uc.db.Faucet().Update(ctx, faucetModel, userInfo.AllowedTenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(updateFaucetComponent)
	}

	faucetRetrieved, err := uc.db.Faucet().FindOneByUUID(ctx, faucet.UUID, userInfo.AllowedTenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(updateFaucetComponent)
	}

	logger.Info("faucet updated successfully")
	return parsers.NewFaucetFromModel(faucetRetrieved), nil
}
