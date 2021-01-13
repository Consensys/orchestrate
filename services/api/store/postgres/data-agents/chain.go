package dataagents

import (
	"context"

	gopg "github.com/go-pg/pg/v9"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models"

	"github.com/gofrs/uuid"
	pg "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
)

const chainDAComponent = "data-agents.chain"

// PGChain is a Chain data agent for PostgreSQL
type PGChain struct {
	db pg.DB
}

// NewPGChain creates a new PGChain
func NewPGChain(db pg.DB) store.ChainAgent {
	return &PGChain{db: db}
}

// Insert Inserts a new chain in DB
func (agent *PGChain) Insert(ctx context.Context, chain *models.Chain) error {
	if chain.UUID == "" {
		chain.UUID = uuid.Must(uuid.NewV4()).String()
	}

	err := pg.Insert(ctx, agent.db, chain)
	if err != nil {
		return errors.FromError(err).ExtendComponent(chainDAComponent)
	}

	return nil
}

// FindOneByUUID Finds a chain in DB by UUID
func (agent *PGChain) FindOneByUUID(ctx context.Context, chainUUID string, tenants []string) (*models.Chain, error) {
	chain := &models.Chain{}

	query := agent.db.ModelContext(ctx, chain).Where("uuid = ?", chainUUID).Relation("PrivateTxManagers")
	query = pg.WhereAllowedTenants(query, "tenant_id", tenants)

	err := pg.SelectOne(ctx, query)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(chainDAComponent)
	}

	return chain, nil
}

// FindOneByName Finds a chain in DB by name
func (agent *PGChain) FindOneByName(ctx context.Context, name string, tenants []string) (*models.Chain, error) {
	chain := &models.Chain{}

	query := agent.db.ModelContext(ctx, chain).Where("name = ?", name)
	query = pg.WhereAllowedTenants(query, "tenant_id", tenants)

	err := pg.SelectOne(ctx, query)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(chainDAComponent)
	}

	return chain, nil
}

func (agent *PGChain) Search(ctx context.Context, filters *entities.ChainFilters, tenants []string) ([]*models.Chain, error) {
	var chains []*models.Chain

	query := agent.db.ModelContext(ctx, &chains).Relation("PrivateTxManagers")

	if len(filters.Names) > 0 {
		query = query.Where("name in (?)", gopg.In(filters.Names))
	}

	query = pg.WhereAllowedTenants(query, "tenant_id", tenants).Order("created_at ASC")

	err := pg.Select(ctx, query)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(chainDAComponent)
	}

	return chains, nil
}

func (agent *PGChain) Update(ctx context.Context, chain *models.Chain, tenants []string) error {
	query := agent.db.ModelContext(ctx, chain).Where("uuid = ?", chain.UUID)
	query = pg.WhereAllowedTenantsDefault(query, tenants)

	err := pg.UpdateNotZero(ctx, query)
	if err != nil {
		return errors.FromError(err).ExtendComponent(chainDAComponent)
	}

	return nil
}

func (agent *PGChain) Delete(ctx context.Context, chain *models.Chain, tenants []string) error {
	query := agent.db.ModelContext(ctx, chain).Where("uuid = ?", chain.UUID)
	query = pg.WhereAllowedTenantsDefault(query, tenants)

	err := pg.Delete(ctx, query)
	if err != nil {
		return errors.FromError(err).ExtendComponent(chainDAComponent)
	}

	return nil
}
