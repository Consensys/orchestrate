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

const faucetDAComponent = "data-agents.faucet"

// PGFaucet is a Faucet data agent for PostgreSQL
type PGFaucet struct {
	db pg.DB
}

// NewPGFaucet creates a new PGFaucet
func NewPGFaucet(db pg.DB) store.FaucetAgent {
	return &PGFaucet{db: db}
}

// Insert Inserts a new faucet in DB
func (agent *PGFaucet) Insert(ctx context.Context, faucet *models.Faucet) error {
	if faucet.UUID == "" {
		faucet.UUID = uuid.Must(uuid.NewV4()).String()
	}

	err := pg.Insert(ctx, agent.db, faucet)
	if err != nil {
		return errors.FromError(err).ExtendComponent(faucetDAComponent)
	}

	return nil
}

// FindOneByUUID Finds a faucet in DB
func (agent *PGFaucet) FindOneByUUID(ctx context.Context, faucetUUID string, tenants []string) (*models.Faucet, error) {
	faucet := &models.Faucet{}

	query := agent.db.ModelContext(ctx, faucet).Where("uuid = ?", faucetUUID)
	query = pg.WhereAllowedTenants(query, "tenant_id", tenants)

	err := pg.SelectOne(ctx, query)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(faucetDAComponent)
	}

	return faucet, nil
}

func (agent *PGFaucet) Search(ctx context.Context, filters *entities.FaucetFilters, tenants []string) ([]*models.Faucet, error) {
	var faucets []*models.Faucet

	query := agent.db.ModelContext(ctx, &faucets)
	if len(filters.Names) > 0 {
		query = query.Where("name in (?)", gopg.In(filters.Names))
	}

	if filters.ChainRule != "" {
		query = query.Where("chain_rule = ?", filters.ChainRule)
	}

	query = pg.WhereAllowedTenants(query, "tenant_id", tenants).Order("created_at ASC")

	err := pg.Select(ctx, query)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(faucetDAComponent)
	}

	return faucets, nil
}

func (agent *PGFaucet) Update(ctx context.Context, faucet *models.Faucet, tenants []string) error {
	query := agent.db.ModelContext(ctx, faucet).Where("uuid = ?", faucet.UUID)
	query = pg.WhereAllowedTenantsDefault(query, tenants)

	err := pg.UpdateNotZero(ctx, query)
	if err != nil {
		return errors.FromError(err).ExtendComponent(faucetDAComponent)
	}

	return nil
}

func (agent *PGFaucet) Delete(ctx context.Context, faucet *models.Faucet, tenants []string) error {
	query := agent.db.ModelContext(ctx, faucet).Where("uuid = ?", faucet.UUID)
	query = pg.WhereAllowedTenantsDefault(query, tenants)

	err := pg.Delete(ctx, query)
	if err != nil {
		return errors.FromError(err).ExtendComponent(faucetDAComponent)
	}

	return nil
}
