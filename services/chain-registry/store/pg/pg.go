package pg

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v9"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

// ChainRegistry is a traefik dynamic config registry based on PostgreSQL
type ChainRegistry struct {
	db *pg.DB
}

// NewChainRegistry creates a new chain registry
func NewChainRegistry(db *pg.DB) *ChainRegistry {
	return &ChainRegistry{db: db}
}

// NewContractRegistryFromPGOptions creates a new pg chain registry
func NewChainRegistryPGOptions(opts *pg.Options) *ChainRegistry {
	return NewChainRegistry(pg.Connect(opts))
}

func (r *ChainRegistry) RegisterChain(ctx context.Context, chain *types.Chain) error {
	chain.SetDefault()

	_, err := r.db.ModelContext(ctx, chain).
		Insert()
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	return nil
}

func (r *ChainRegistry) GetChains(ctx context.Context, filters map[string]string) ([]*types.Chain, error) {
	chains := make([]*types.Chain, 0)

	req := r.db.ModelContext(ctx, &chains)
	for k, v := range filters {
		req.Where(fmt.Sprintf("%s = ?", k), v)
	}

	err := req.Select()
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return chains, nil
}

func (r *ChainRegistry) GetChainsByTenant(ctx context.Context, filters map[string]string, tenantID string) ([]*types.Chain, error) {
	chains := make([]*types.Chain, 0)

	req := r.db.ModelContext(ctx, &chains).
		Where("tenant_id = ?", tenantID)
	for k, v := range filters {
		req.Where(fmt.Sprintf("%s = ?", k), v)
	}

	err := req.Select()
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return chains, nil
}

func (r *ChainRegistry) GetChainByUUID(ctx context.Context, uuid string) (*types.Chain, error) {
	chain := &types.Chain{}

	err := r.db.ModelContext(ctx, chain).Where("uuid = ?", uuid).Select()
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return chain, nil
}

func (r *ChainRegistry) GetChainByUUIDAndTenant(ctx context.Context, uuid, tenantID string) (*types.Chain, error) {
	chain := &types.Chain{}

	err := r.db.ModelContext(ctx, chain).Where("uuid = ?", uuid).Where("tenant_id = ?", tenantID).Select()
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return chain, nil
}

func (r *ChainRegistry) UpdateChainByName(ctx context.Context, chain *types.Chain) error {
	res, err := r.db.ModelContext(ctx, chain).Where("name = ?", chain.Name).UpdateNotZero()
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		return errors.NotFoundError("no chain found with tenant_id=%s and name=%s", chain.TenantID, chain.Name).ExtendComponent(component)
	}
	return nil
}

func (r *ChainRegistry) UpdateChainByUUID(ctx context.Context, chain *types.Chain) error {
	res, err := r.db.ModelContext(ctx, chain).WherePK().UpdateNotZero()
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		return errors.NotFoundError("no chain found with uuid %s and tenant_id %s", chain.UUID, chain.TenantID).ExtendComponent(component)
	}
	return nil
}

func (r *ChainRegistry) DeleteChainByUUID(ctx context.Context, uuid string) error {
	chain := &types.Chain{}

	res, err := r.db.ModelContext(ctx, chain).Where("uuid = ?", uuid).Delete()
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		return errors.NotFoundError("no chain found with uuid %s", uuid).ExtendComponent(component)
	}

	return nil
}

func (r *ChainRegistry) DeleteChainByUUIDAndTenant(ctx context.Context, uuid, tenantID string) error {
	chain := &types.Chain{}

	res, err := r.db.ModelContext(ctx, chain).Where("uuid = ?", uuid).Where("tenant_id = ?", tenantID).Delete()
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		return errors.NotFoundError("no chain found with uuid %s and tenant_id %s", uuid, tenantID).ExtendComponent(component)
	}

	return nil
}
