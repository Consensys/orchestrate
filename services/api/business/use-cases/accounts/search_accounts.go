package accounts

import (
	"context"

	parsers2 "github.com/ConsenSys/orchestrate/services/api/business/parsers"
	usecases "github.com/ConsenSys/orchestrate/services/api/business/use-cases"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/log"
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/ConsenSys/orchestrate/services/api/store"
)

const searchAccountsComponent = "use-cases.search-accounts"

type searchAccountsUseCase struct {
	db     store.DB
	logger *log.Logger
}

func NewSearchAccountsUseCase(db store.DB) usecases.SearchAccountsUseCase {
	return &searchAccountsUseCase{
		db:     db,
		logger: log.NewLogger().SetComponent(searchAccountsComponent),
	}
}

func (uc *searchAccountsUseCase) Execute(ctx context.Context, filters *entities.AccountFilters, tenants []string) ([]*entities.Account, error) {
	models, err := uc.db.Account().Search(ctx, filters, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(searchAccountsComponent)
	}

	var resp []*entities.Account
	for _, model := range models {
		iden := parsers2.NewAccountEntityFromModels(model)
		resp = append(resp, iden)
	}

	uc.logger.WithContext(ctx).Debug("accounts found successfully")
	return resp, nil
}
