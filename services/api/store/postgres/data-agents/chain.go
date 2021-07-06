package dataagents

import (
	"context"
	"time"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
	pg "github.com/ConsenSys/orchestrate/pkg/toolkit/database/postgres"
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/ConsenSys/orchestrate/services/api/store"
	"github.com/ConsenSys/orchestrate/services/api/store/models"
	gopg "github.com/go-pg/pg/v9"
	"github.com/gofrs/uuid"
)

const chainDAComponent = "data-agents.chain"

// PGChain is a Chain data agent for PostgreSQL
type PGChain struct {
	db     pg.DB
	logger *log.Logger
}

// NewPGChain creates a new PGChain
func NewPGChain(db pg.DB) store.ChainAgent {
	return &PGChain{db: db, logger: log.NewLogger().SetComponent(chainDAComponent)}
}

// Insert Inserts a new chain in DB
func (agent *PGChain) Insert(ctx context.Context, chain *models.Chain) error {
	if chain.UUID == "" {
		chain.UUID = uuid.Must(uuid.NewV4()).String()
	}

	err := pg.Insert(ctx, agent.db, chain)
	if err != nil {
		agent.logger.WithContext(ctx).WithError(err).Error("failed to insert chain")
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
		if !errors.IsNotFoundError(err) {
			agent.logger.WithContext(ctx).WithError(err).Error("failed to select chain by uuid")
		}
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
		if !errors.IsNotFoundError(err) {
			agent.logger.WithContext(ctx).WithError(err).Error("failed to select chain by name")
		}
		return nil, errors.FromError(err).ExtendComponent(chainDAComponent)
	}

	return chain, nil
}

func (agent *PGChain) Search(ctx context.Context, filters *entities.ChainFilters, tenants []string) ([]*models.Chain, error) {
	var chains []*models.Chain

	query := agent.db.ModelContext(ctx, &chains)

	if len(filters.Names) > 0 {
		query = query.Where("name in (?)", gopg.In(filters.Names))
	}

	query = pg.WhereAllowedTenants(query, "tenant_id", tenants).Order("created_at ASC")

	if err := pg.Select(ctx, query); err != nil {
		if !errors.IsNotFoundError(err) {
			agent.logger.WithContext(ctx).WithError(err).Error("failed to search chains")
		}
		return nil, errors.FromError(err).ExtendComponent(chainDAComponent)
	}

	if len(chains) == 0 {
		return chains, nil
	}

	// We manually link chains to privateTxManager avoiding GO-PG strategy of one query per chain
	chainUUIDs := make([]string, len(chains))
	chainsMap := map[string]*models.Chain{}
	for idx, c := range chains {
		chainUUIDs[idx] = c.UUID
		chainsMap[c.UUID] = c
	}

	var privateTxManagers []*models.PrivateTxManager
	query = agent.db.ModelContext(ctx, &privateTxManagers)
	query = query.Where("chain_uuid in (?)", gopg.In(chainUUIDs))

	if err := pg.Select(ctx, query); err != nil {
		agent.logger.WithContext(ctx).WithError(err).Error("failed to search for private tx manager")
		return nil, errors.FromError(err).ExtendComponent(chainDAComponent)
	}

	for _, pm := range privateTxManagers {
		chainsMap[pm.ChainUUID].PrivateTxManagers = []*models.PrivateTxManager{pm}
	}

	return chains, nil
}

func (agent *PGChain) Update(ctx context.Context, chain *models.Chain, tenants []string) error {
	chain.UpdatedAt = time.Now().UTC()
	query := agent.db.ModelContext(ctx, chain).Where("uuid = ?", chain.UUID)
	query = pg.WhereAllowedTenantsDefault(query, tenants)

	err := pg.UpdateNotZero(ctx, query)
	if err != nil {
		agent.logger.WithContext(ctx).WithError(err).Error("failed to update chain")
		return errors.FromError(err).ExtendComponent(chainDAComponent)
	}

	return nil
}

func (agent *PGChain) Delete(ctx context.Context, chain *models.Chain, tenants []string) error {
	query := agent.db.ModelContext(ctx, chain).Where("uuid = ?", chain.UUID)
	query = pg.WhereAllowedTenantsDefault(query, tenants)

	err := pg.Delete(ctx, query)
	if err != nil {
		agent.logger.WithContext(ctx).WithError(err).Error("failed to deleted chain")
		return errors.FromError(err).ExtendComponent(chainDAComponent)
	}

	return nil
}
