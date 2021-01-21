package accounts

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/parsers"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
)

const getAccountComponent = "use-cases.get-account"

type getAccountUseCase struct {
	db     store.DB
	logger *log.Logger
}

func NewGetAccountUseCase(db store.DB) usecases.GetAccountUseCase {
	return &getAccountUseCase{
		db:     db,
		logger: log.NewLogger().SetComponent(getAccountComponent),
	}
}

func (uc *getAccountUseCase) Execute(ctx context.Context, address string, tenants []string) (*entities.Account, error) {
	ctx = log.WithFields(ctx, log.Field("address", address))
	logger := uc.logger.WithContext(ctx)

	model, err := uc.db.Account().FindOneByAddress(ctx, address, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getAccountComponent)
	}

	logger.Debug("account found successfully")
	return parsers.NewAccountEntityFromModels(model), nil
}
