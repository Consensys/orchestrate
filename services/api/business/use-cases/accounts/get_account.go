package accounts

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/parsers"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
)

const getAccountComponent = "use-cases.get-account"

type getAccountUseCase struct {
	db store.DB
}

func NewGetAccountUseCase(db store.DB) usecases.GetAccountUseCase {
	return &getAccountUseCase{
		db: db,
	}
}

func (uc *getAccountUseCase) Execute(ctx context.Context, address string, tenants []string) (*entities.Account, error) {
	logger := log.WithContext(ctx).WithField("address", address).WithField("tenants", tenants)
	logger.Debug("getting account")

	model, err := uc.db.Account().FindOneByAddress(ctx, address, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getAccountComponent)
	}

	log.WithContext(ctx).Debug("account found successfully")
	return parsers.NewAccountEntityFromModels(model), nil
}
