package types

import "context"

type ChainRegistryStore interface {
	RegisterConfig(ctx context.Context, config *Config) error
	RegisterConfigs(ctx context.Context, configs *[]Config) error
	GetConfig(ctx context.Context) ([]*Config, error)
	GetConfigByID(ctx context.Context, config *Config) error
	GetConfigByTenantID(ctx context.Context, config *Config) ([]*Config, error)
	UpdateConfigByID(ctx context.Context, config *Config) error
	DeregisterConfigByID(ctx context.Context, config *Config) error
	DeregisterConfigsByIds(ctx context.Context, config *[]Config) error
	DeregisterConfigByTenantID(ctx context.Context, config *Config) error
}
