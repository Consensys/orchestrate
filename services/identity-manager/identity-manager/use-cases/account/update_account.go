package account

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/identity-manager/parsers"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/identity-manager/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/store"
)

const updateAccountComponent = "use-cases.update-account"

type updateAccountUseCase struct {
	db store.DB
}

func NewUpdateAccountUseCase(db store.DB) usecases.UpdateAccountUseCase {
	return &updateAccountUseCase{
		db: db,
	}
}

func (uc *updateAccountUseCase) Execute(ctx context.Context, account *entities.Account, tenants []string) (*entities.Account, error) {
	log.WithContext(ctx).WithField("address", account.Address).
		WithField("tenants", tenants).
		Debug("updating account")

	model, err := uc.db.Account().FindOneByAddress(ctx, account.Address, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(updateAccountComponent)
	}

	if account.Attributes != nil {
		model.Attributes = account.Attributes
	}
	if account.Alias != "" {
		model.Alias = account.Alias
	}

	err = uc.db.Account().Update(ctx, model)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(updateAccountComponent)
	}

	resp := parsers.NewAccountEntityFromModels(model)

	log.WithContext(ctx).Debug("account updated successfully")
	return resp, nil
}
