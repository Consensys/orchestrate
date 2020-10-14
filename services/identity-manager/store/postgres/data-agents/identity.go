package dataagents

import (
	"context"

	gopg "github.com/go-pg/pg/v9"
	pg "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/store/models"
)

const identityDAComponent = "data-agents.identity"

// NewPGIdentity creates a new PGIdentity
func NewPGIdentity(db pg.DB) *PGIdentity {
	return &PGIdentity{db: db}
}

// PGIdentity is an Identity data agent for PostgreSQL
type PGIdentity struct {
	db pg.DB
}

func (agent *PGIdentity) Insert(ctx context.Context, identity *models.Identity) error {
	agent.db.ModelContext(ctx, identity)
	err := pg.Insert(ctx, agent.db, identity)
	if err != nil {
		return errors.FromError(err).ExtendComponent(identityDAComponent)
	}

	return nil
}

func (agent *PGIdentity) Search(ctx context.Context, filters *entities.IdentityFilters, tenants []string) ([]*models.Identity, error) {
	var idens []*models.Identity

	query := agent.db.ModelContext(ctx, &idens)
	if len(filters.Aliases) > 0 {
		query = query.Where("alias in (?)", gopg.In(filters.Aliases))
	}

	query = pg.WhereAllowedTenants(query, "tenant_id", tenants).
		Order("id ASC")

	err := pg.Select(ctx, query)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(identityDAComponent)
	}

	return idens, nil
}
