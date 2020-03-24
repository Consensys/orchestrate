package pg

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

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

func (r *ChainRegistry) RegisterChain(ctx context.Context, chain *types.Chain) error {
	_, err := r.db.ModelContext(ctx, chain).Insert()
	if err != nil {
		if errors.IsAlreadyExistsError(err) {
			errMessage := "Chain already exists"
			log.WithError(err).Error(errMessage)
			return errors.AlreadyExistsError(errMessage).ExtendComponent(component)
		}

		errMessage := "Failed to register chain"
		log.WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	return nil
}

func (r *ChainRegistry) GetChains(ctx context.Context, filters map[string]string) ([]*types.Chain, error) {
	var chains []*types.Chain

	req := r.db.ModelContext(ctx, &chains)
	for k, v := range filters {
		req.Where(fmt.Sprintf("%s = ?", k), v)
	}

	err := req.Select()
	if err != nil {
		errMessage := "Failed to get chains"
		log.WithError(err).Error(errMessage)
		return nil, errors.PostgresConnectionError(errMessage).ExtendComponent(component)
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
		errMessage := "Failed to get chains for tenant ID %s"
		log.WithError(err).Error(errMessage, tenantID)
		return nil, errors.PostgresConnectionError(errMessage, tenantID).ExtendComponent(component)
	}

	return chains, nil
}

func (r *ChainRegistry) GetChainByUUID(ctx context.Context, chainUUID string) (*types.Chain, error) {
	chain := &types.Chain{}

	err := r.db.ModelContext(ctx, chain).Where("uuid = ?", chainUUID).Select()
	if err != nil && err == pg.ErrNoRows {
		errMessage := "could not load chain with chainUUID: %s"
		log.WithError(err).Debugf(errMessage, chainUUID)
		return nil, errors.NotFoundError(errMessage, chainUUID).ExtendComponent(component)
	} else if err != nil {
		errMessage := "Failed to get chain by UUID"
		log.WithError(err).Error(errMessage)
		return nil, errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	return chain, nil
}

func (r *ChainRegistry) GetChainByUUIDAndTenant(ctx context.Context, chainUUID, tenantID string) (*types.Chain, error) {
	chain := &types.Chain{}

	err := r.db.ModelContext(ctx, chain).Where("uuid = ?", chainUUID).Where("tenant_id = ?", tenantID).Select()
	if err != nil && err == pg.ErrNoRows {
		errMessage := "could not load chain with chainUUID: %s and tenant: %s"
		log.WithError(err).Debugf(errMessage, chainUUID, tenantID)
		return nil, errors.NotFoundError(errMessage, chainUUID, tenantID).ExtendComponent(component)
	} else if err != nil {
		errMessage := "Failed to get chain by UUID and tenant"
		log.WithError(err).Error(errMessage)
		return nil, errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	return chain, nil
}

func (r *ChainRegistry) UpdateChainByName(ctx context.Context, chain *types.Chain) error {
	res, err := r.db.ModelContext(ctx, chain).Where("name = ?", chain.Name).UpdateNotZero()
	if err != nil {
		errMessage := "Failed to update chain by name"
		log.WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		errMessage := "no chain found with tenant_id=%s and name=%s"
		log.WithError(err).Error(errMessage, chain.TenantID, chain.Name)
		return errors.NotFoundError(errMessage, chain.TenantID, chain.Name).ExtendComponent(component)
	}

	return nil
}

func (r *ChainRegistry) UpdateChainByUUID(ctx context.Context, chain *types.Chain) error {
	res, err := r.db.ModelContext(ctx, chain).WherePK().UpdateNotZero()
	if err != nil {
		errMessage := "Failed to update chain by UUID"
		log.WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		errMessage := "no chain found with uuid %s for update"
		log.WithError(err).Error(errMessage, chain.UUID)
		return errors.NotFoundError(errMessage, chain.UUID).ExtendComponent(component)
	}

	return nil
}

func (r *ChainRegistry) DeleteChainByUUID(ctx context.Context, chainUUID string) error {
	chain := &types.Chain{}

	res, err := r.db.ModelContext(ctx, chain).Where("uuid = ?", chainUUID).Delete()
	if err != nil {
		errMessage := "Failed to delete chain by UUID"
		log.WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		errMessage := "no chain found with uuid %s for delete"
		log.WithError(err).Error(errMessage, chainUUID)
		return errors.NotFoundError(errMessage, chainUUID).ExtendComponent(component)
	}

	return nil
}

func (r *ChainRegistry) DeleteChainByUUIDAndTenant(ctx context.Context, chainUUID, tenantID string) error {
	chain := &types.Chain{}

	res, err := r.db.ModelContext(ctx, chain).Where("uuid = ?", chainUUID).Where("tenant_id = ?", tenantID).Delete()
	if err != nil {
		errMessage := "Failed to delete chain by UUID and tenant"
		log.WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		errMessage := "no chain found with uuid %s and tenant_id %s"
		log.WithError(err).Error(errMessage, chainUUID, tenantID)
		return errors.NotFoundError(errMessage, chainUUID, tenantID).ExtendComponent(component)
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
		errMessage := "Failed to register faucet"
		log.WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	return nil
}

func (r *ChainRegistry) UpdateFaucetByUUID(ctx context.Context, faucet *types.Faucet) error {
	res, err := r.db.ModelContext(ctx, faucet).
		WherePK().
		UpdateNotZero()
	if err != nil {
		errMessage := "Failed to update faucet"
		log.WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		errMessage := "no faucet found with uuid %s"
		log.WithError(err).Error(errMessage, faucet.UUID)
		return errors.NotFoundError(errMessage, faucet.UUID).ExtendComponent(component)
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
		errMessage := "Failed to get faucets"
		log.WithError(err).Error(errMessage)
		return nil, errors.PostgresConnectionError(errMessage).ExtendComponent(component)
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
		errMessage := "Failed to get faucets for tenant ID %s"
		log.WithError(err).Error(errMessage, tenantID)
		return nil, errors.PostgresConnectionError(errMessage, tenantID).ExtendComponent(component)
	}

	return faucets, nil
}

func (r *ChainRegistry) GetFaucetByUUID(ctx context.Context, chainUUID string) (*types.Faucet, error) {
	faucet := &types.Faucet{}

	err := r.db.ModelContext(ctx, faucet).Where("uuid = ?", chainUUID).Select()
	if err != nil && err == pg.ErrNoRows {
		errMessage := "could not load faucet with chainUUID: %s"
		log.WithError(err).Debugf(errMessage, chainUUID)
		return nil, errors.NotFoundError(errMessage, chainUUID).ExtendComponent(component)
	} else if err != nil {
		errMessage := "Failed to get faucet"
		log.WithError(err).Error(errMessage)
		return nil, errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	return faucet, nil
}

func (r *ChainRegistry) GetFaucetByUUIDAndTenant(ctx context.Context, chainUUID, tenantID string) (*types.Faucet, error) {
	faucet := &types.Faucet{}

	err := r.db.ModelContext(ctx, faucet).Where("uuid = ?", chainUUID).Where("tenant_id = ?", tenantID).Select()
	if err != nil && err == pg.ErrNoRows {
		errMessage := "could not load faucet with chainUUID: %s and tenant: %s"
		log.WithError(err).Debugf(errMessage, chainUUID, tenantID)
		return nil, errors.NotFoundError(errMessage, chainUUID, tenantID).ExtendComponent(component)
	} else if err != nil {
		errMessage := "Failed to get faucet from DB"
		log.WithError(err).Error(errMessage)
		return nil, errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	return faucet, nil
}

func (r *ChainRegistry) DeleteFaucetByUUID(ctx context.Context, chainUUID string) error {
	faucet := &types.Faucet{}

	res, err := r.db.ModelContext(ctx, faucet).Where("uuid = ?", chainUUID).Delete()
	if err != nil {
		errMessage := "Failed to delete faucet by UUID"
		log.WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		errMessage := "no faucet found with chainUUID: %s"
		log.WithError(err).Error(errMessage, chainUUID)
		return errors.NotFoundError(errMessage, chainUUID).ExtendComponent(component)
	}

	return nil
}

func (r *ChainRegistry) DeleteFaucetByUUIDAndTenant(ctx context.Context, chainUUID, tenantID string) error {
	faucet := &types.Faucet{}

	res, err := r.db.ModelContext(ctx, faucet).Where("uuid = ?", chainUUID).Where("tenant_id = ?", tenantID).Delete()
	if err != nil {
		errMessage := "Failed to delete faucet by UUID and tenant"
		log.WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		errMessage := "no faucet found with uuid %s and tenant_id %s"
		log.WithError(err).Error(errMessage, chainUUID, tenantID)
		return errors.NotFoundError(errMessage, chainUUID, tenantID).ExtendComponent(component)
	}

	return nil
}
