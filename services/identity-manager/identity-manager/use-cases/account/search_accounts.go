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

const searchAccountsComponent = "use-cases.search-accounts"

type searchAccountsUseCase struct {
	db store.DB
}

func NewSearchAccountsUseCase(db store.DB) usecases.SearchAccountsUseCase {
	return &searchAccountsUseCase{
		db: db,
	}
}

func (uc *searchAccountsUseCase) Execute(ctx context.Context, filters *entities.AccountFilters, tenants []string) ([]*entities.Account, error) {
	log.WithContext(ctx).WithField("filters", filters).
		WithField("tenants", tenants).
		Debug("searching accounts")

	models, err := uc.db.Account().Search(ctx, filters, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(searchAccountsComponent)
	}

	var resp []*entities.Account
	for _, model := range models {
		iden := parsers.NewAccountEntityFromModels(model)
		resp = append(resp, iden)
	}

	log.WithContext(ctx).Debug("accounts found successfully")
	return resp, nil
}
