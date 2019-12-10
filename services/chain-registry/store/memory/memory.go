package memory

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

type ChainRegistry struct {
}

// NewChainRegistry creates a new chain registry
func NewChainRegistry() *ChainRegistry {
	return &ChainRegistry{}
}

// TODO: Implement in memory registry
func (r *ChainRegistry) RegisterConfig(ctx context.Context, config *types.Config) error {
	return nil
}

func (r *ChainRegistry) RegisterConfigs(ctx context.Context, configs *[]types.Config) error {
	return nil
}

func (r *ChainRegistry) GetConfig(ctx context.Context) ([]*types.Config, error) {
	return nil, nil
}

func (r *ChainRegistry) GetConfigByID(ctx context.Context, config *types.Config) error {
	return nil
}

func (r *ChainRegistry) GetConfigByTenantID(ctx context.Context, config *types.Config) ([]*types.Config, error) {
	return nil, nil
}

func (r *ChainRegistry) UpdateConfigByID(ctx context.Context, config *types.Config) error {
	return nil
}

func (r *ChainRegistry) DeregisterConfigByID(ctx context.Context, config *types.Config) error {
	return nil
}

func (r *ChainRegistry) DeregisterConfigsByIds(ctx context.Context, configs *[]types.Config) error {
	return nil
}

func (r *ChainRegistry) DeregisterConfigByTenantID(ctx context.Context, config *types.Config) error {
	return nil
}
