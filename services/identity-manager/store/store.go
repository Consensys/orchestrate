package store

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/identity-manager/store/models"
)

//go:generate mockgen -source=store.go -destination=mocks/mock.go -package=mocks

type Agents interface {
	Account() AccountAgent
}

type DB interface {
	database.DB
	Agents
}

// Interfaces data agents
type AccountAgent interface {
	Insert(ctx context.Context, identity *models.Account) error
	Update(ctx context.Context, identity *models.Account) error
	FindOneByAddress(ctx context.Context, address string, tenants []string) (*models.Account, error)
	Search(ctx context.Context, filters *entities.AccountFilters, tenants []string) ([]*models.Account, error)
}
