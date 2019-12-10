package pg

import (
	"context"

	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
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

// RegisterConfig register a config
func (r *ChainRegistry) RegisterConfig(ctx context.Context, config *types.Config) error {

	configStruct, err := types.GetConfigStruct(config.ConfigType)
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	err = types.UnmarshalJSONConfig(config.Config, configStruct)
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	_, err = r.db.ModelContext(ctx, config).
		Insert()
	if err != nil {
		log.WithError(err).Debugf("could not register configs")
		return errors.FromError(err).ExtendComponent(component)
	}

	return nil
}

// RegisterConfigs register many configs at once
func (r *ChainRegistry) RegisterConfigs(ctx context.Context, configs *[]types.Config) error {

	for _, config := range *configs {
		configStruct, err := types.GetConfigStruct(config.ConfigType)
		if err != nil {
			return errors.FromError(err).ExtendComponent(component)
		}

		err = types.UnmarshalJSONConfig(config.Config, configStruct)
		if err != nil {
			return errors.FromError(err).ExtendComponent(component)
		}
	}

	_, err := r.db.ModelContext(ctx, configs).
		Insert()
	if err != nil {
		log.WithError(err).Debugf("could not register configs")
		return errors.FromError(err).ExtendComponent(component)
	}

	return nil
}

// GetConfigById retrieves config with an id
func (r *ChainRegistry) GetConfig(ctx context.Context) ([]*types.Config, error) {
	var configs []*types.Config

	err := r.db.ModelContext(ctx, &configs).
		Select()
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return configs, nil
}

// GetConfigById retrieves config with an id
func (r *ChainRegistry) GetConfigByID(ctx context.Context, config *types.Config) error {

	err := r.db.ModelContext(ctx, config).
		Where("id = ?", config.ID).
		Select()
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	return nil
}

// GetConfigByTenantID retrieves configs of a tenantId
func (r *ChainRegistry) GetConfigByTenantID(ctx context.Context, config *types.Config) ([]*types.Config, error) {
	var configs []*types.Config

	err := r.db.ModelContext(ctx, &configs).
		Where("tenant_id = ?", config.TenantID).
		Select()
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return configs, nil
}

// UpdateConfigByID updates a config
func (r *ChainRegistry) UpdateConfigByID(ctx context.Context, config *types.Config) error {
	// Check if config match struct
	configStruct, err := types.GetConfigStruct(config.ConfigType)
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	err = types.UnmarshalJSONConfig(config.Config, configStruct)
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	_, err = r.db.ModelContext(ctx, config).
		Set("name = ?name").
		Set("config_type = ?config_type").
		Set("config = ?config").
		Where("id = ?id").
		Update()
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	return nil
}

// DeregisterConfigByID deletes a config
func (r *ChainRegistry) DeregisterConfigByID(ctx context.Context, config *types.Config) error {
	_, err := r.db.ModelContext(ctx, config).
		Where("id = ?id").
		Delete()
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	return nil
}

// DeregisterConfigsByIds deletes many configs
func (r *ChainRegistry) DeregisterConfigsByIds(ctx context.Context, configs *[]types.Config) error {
	_, err := r.db.ModelContext(ctx, configs).
		WherePK().
		Delete()
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	return nil
}

// DeregisterConfigByTenantID deletes configs of a tenant
func (r *ChainRegistry) DeregisterConfigByTenantID(ctx context.Context, config *types.Config) error {
	_, err := r.db.ModelContext(ctx, config).
		Where("tenant_id = ?tenant_id").
		Delete()
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	return nil
}
