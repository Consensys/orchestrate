package account

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/identity-manager/identity-manager/parsers"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/identity-manager/identity-manager/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/identity-manager/store"
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
	log.WithContext(ctx).WithField("address", address).
		WithField("tenants", tenants).
		Debug("getting accounts")

	model, err := uc.db.Account().FindOneByAddress(ctx, address, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getAccountComponent)
	}

	resp := parsers.NewAccountEntityFromModels(model)

	log.WithContext(ctx).Debug("account found successfully")
	return resp, nil
}
