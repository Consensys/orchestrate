package store

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/store/models"
)

//go:generate mockgen -source=store.go -destination=mocks/mock.go -package=mocks

type Agents interface {
	Identity() IdentityAgent
}

type DB interface {
	database.DB
	Agents
}

// Interfaces data agents
type IdentityAgent interface {
	Insert(ctx context.Context, identity *models.Identity) error
	Search(ctx context.Context, filters *entities.IdentityFilters, tenants []string) ([]*models.Identity, error)
}
