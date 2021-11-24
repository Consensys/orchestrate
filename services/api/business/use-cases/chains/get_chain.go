package chains

import (
	"context"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/services/api/business/parsers"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
	"github.com/consensys/orchestrate/services/api/store"
)

const getChainComponent = "use-cases.get-chain"

// getChainUseCase is a use case to get a faucet
type getChainUseCase struct {
	db     store.DB
	logger *log.Logger
}

// NewGetChainUseCase creates a new GetChainUseCase
func NewGetChainUseCase(db store.DB) usecases.GetChainUseCase {
	return &getChainUseCase{
		db:     db,
		logger: log.NewLogger().SetComponent(getChainComponent),
	}
}

// Execute gets a chain
func (uc *getChainUseCase) Execute(ctx context.Context, uuid string, userInfo *multitenancy.UserInfo) (*entities.Chain, error) {
	ctx = log.WithFields(ctx, log.Field("chain", uuid))
	logger := uc.logger.WithContext(ctx)

	chainModel, err := uc.db.Chain().FindOneByUUID(ctx, uuid, userInfo.AllowedTenants, userInfo.Username)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getChainComponent)
	}

	logger.Debug("chain found successfully")
	return parsers.NewChainFromModel(chainModel), nil
}
