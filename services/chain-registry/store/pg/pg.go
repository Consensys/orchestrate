package pg

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"

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
	_, err := r.db.ModelContext(ctx, chain).Insert()
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

func (r *ChainRegistry) GetChainByUUID(ctx context.Context, chainUUID string) (*types.Chain, error) {
	chain := &types.Chain{}

	err := r.db.ModelContext(ctx, chain).Where("uuid = ?", chainUUID).Select()
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return chain, nil
}

func (r *ChainRegistry) GetChainByUUIDAndTenant(ctx context.Context, chainUUID, tenantID string) (*types.Chain, error) {
	chain := &types.Chain{}

	err := r.db.ModelContext(ctx, chain).Where("uuid = ?", chainUUID).Where("tenant_id = ?", tenantID).Select()
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
		return errors.NotFoundError("no chain found with uuid %s", chain.UUID).ExtendComponent(component)
	}
	return nil
}

func (r *ChainRegistry) DeleteChainByUUID(ctx context.Context, chainUUID string) error {
	chain := &types.Chain{}

	res, err := r.db.ModelContext(ctx, chain).Where("uuid = ?", chainUUID).Delete()
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		return errors.NotFoundError("no chain found with uuid %s", chainUUID).ExtendComponent(component)
	}

	return nil
}

func (r *ChainRegistry) DeleteChainByUUIDAndTenant(ctx context.Context, chainUUID, tenantID string) error {
	chain := &types.Chain{}

	res, err := r.db.ModelContext(ctx, chain).Where("uuid = ?", chainUUID).Where("tenant_id = ?", tenantID).Delete()
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		return errors.NotFoundError("no chain found with uuid %s and tenant_id %s", chainUUID, tenantID).ExtendComponent(component)
	}

	return nil
}

func (r *ChainRegistry) RegisterFaucet(ctx context.Context, faucet *types.Faucet) error {
	if faucet.UUID == "" {
		faucet.UUID = uuid.NewV4().String()
	}
	_, err := r.db.ModelContext(ctx, faucet).
		Insert()
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	return nil
}

func (r *ChainRegistry) UpdateFaucetByUUID(ctx context.Context, faucet *types.Faucet) error {
	res, err := r.db.ModelContext(ctx, faucet).
		WherePK().
		UpdateNotZero()
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		return errors.NotFoundError("no faucet found with uuid %s", faucet.UUID).ExtendComponent(component)
	}
	return nil
}

func (r *ChainRegistry) GetFaucets(ctx context.Context, filters map[string]string) ([]*types.Faucet, error) {
	faucets := make([]*types.Faucet, 0)

	req := r.db.ModelContext(ctx, &faucets)
	for k, v := range filters {
		req.Where(fmt.Sprintf("%s = ?", k), v)
	}

	err := req.Select()
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return faucets, nil
}

func (r *ChainRegistry) GetFaucetsByTenant(ctx context.Context, filters map[string]string, tenantID string) ([]*types.Faucet, error) {
	faucets := make([]*types.Faucet, 0)

	req := r.db.ModelContext(ctx, &faucets).
		Where("tenant_id = ?", tenantID)
	for k, v := range filters {
		req.Where(fmt.Sprintf("%s = ?", k), v)
	}

	err := req.Select()
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return faucets, nil
}

func (r *ChainRegistry) GetFaucetByUUID(ctx context.Context, chainUUID string) (*types.Faucet, error) {
	faucet := &types.Faucet{}

	err := r.db.ModelContext(ctx, faucet).Where("uuid = ?", chainUUID).Select()
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return faucet, nil
}

func (r *ChainRegistry) GetFaucetByUUIDAndTenant(ctx context.Context, chainUUID, tenantID string) (*types.Faucet, error) {
	faucet := &types.Faucet{}

	err := r.db.ModelContext(ctx, faucet).Where("uuid = ?", chainUUID).Where("tenant_id = ?", tenantID).Select()
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return faucet, nil
}

func (r *ChainRegistry) DeleteFaucetByUUID(ctx context.Context, chainUUID string) error {
	faucet := &types.Faucet{}

	res, err := r.db.ModelContext(ctx, faucet).Where("uuid = ?", chainUUID).Delete()
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		return errors.NotFoundError("no faucet found with uuid %s", chainUUID).ExtendComponent(component)
	}

	return nil
}

func (r *ChainRegistry) DeleteFaucetByUUIDAndTenant(ctx context.Context, chainUUID, tenantID string) error {
	faucet := &types.Faucet{}

	res, err := r.db.ModelContext(ctx, faucet).Where("uuid = ?", chainUUID).Where("tenant_id = ?", tenantID).Delete()
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		return errors.NotFoundError("no faucet found with uuid %s and tenant_id %s", chainUUID, tenantID).ExtendComponent(component)
	}

	return nil
}
