package store

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	pgstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/pg"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

const StoreName = "chains"

//go:generate mockgen -source=store.go -destination=mock/mock.go -package=mock

type ChainStore interface {
	RegisterChain(ctx context.Context, chain *types.Chain) error

	GetChains(ctx context.Context, filters map[string]string) ([]*types.Chain, error)
	GetChainsByTenant(ctx context.Context, filters map[string]string, tenantID string) ([]*types.Chain, error)
	GetChainByUUID(ctx context.Context, uuid string) (*types.Chain, error)
	GetChainByUUIDAndTenant(ctx context.Context, uuid string, tenantID string) (*types.Chain, error)

	UpdateChainByName(ctx context.Context, chain *types.Chain) error
	UpdateChainByUUID(ctx context.Context, chain *types.Chain) error

	DeleteChainByUUID(ctx context.Context, uuid string) error
	DeleteChainByUUIDAndTenant(ctx context.Context, uuid string, tenantID string) error
}

type FaucetStore interface {
	RegisterFaucet(ctx context.Context, faucet *types.Faucet) error

	GetFaucets(ctx context.Context, filters map[string]string) ([]*types.Faucet, error)
	GetFaucetsByTenant(ctx context.Context, filters map[string]string, tenantID string) ([]*types.Faucet, error)
	GetFaucetByUUID(ctx context.Context, uuid string) (*types.Faucet, error)
	GetFaucetByUUIDAndTenant(ctx context.Context, uuid string, tenantID string) (*types.Faucet, error)

	UpdateFaucetByUUID(ctx context.Context, faucet *types.Faucet) error

	DeleteFaucetByUUID(ctx context.Context, uuid string) error
	DeleteFaucetByUUIDAndTenant(ctx context.Context, uuid string, tenantID string) error
}

type ChainRegistryStore interface {
	ChainStore
	FaucetStore
}

type Builder interface {
	Build(context.Context, *Config) (ChainRegistryStore, error)
}

func NewBuilder(pgmngr postgres.Manager) Builder {
	return &builder{
		postgres: pgstore.NewBuilder(pgmngr),
	}
}

type builder struct {
	postgres *pgstore.Builder
}

func (b *builder) Build(ctx context.Context, conf *Config) (store ChainRegistryStore, err error) {
	logCtx := log.With(ctx, log.Str("store", StoreName))
	switch conf.Type {
	case postgresType:
		conf.Postgres.PG.ApplicationName = StoreName
		return b.postgres.Build(logCtx, conf.Postgres)
	default:
		return nil, fmt.Errorf("invalid chain registry store type %q", conf.Type)
	}
}

func ImportChains(ctx context.Context, s ChainRegistryStore, chains []string) {
	logger := log.FromContext(ctx)
	for _, v := range chains {
		logger.WithField("config", v).Debugf("import chain from configuration")
		chain := &types.Chain{}
		dec := json.NewDecoder(strings.NewReader(v))
		dec.DisallowUnknownFields() // Force errors if unknown fields
		err := dec.Decode(chain)
		if err != nil {
			logger.WithError(err).Errorf("could not import chain (invalid configuration provided)")
			continue
		}

		chain.SetDefault()

		err = s.RegisterChain(ctx, chain)
		if err != nil {
			logger.WithError(err).Errorf("could not import chain")
			continue
		}

		logger.WithFields(logrus.Fields{
			"chain.name":   chain.Name,
			"chain.uuid":   chain.UUID,
			"chain.tenant": chain.TenantID,
		}).Infof("imported chain from configuration")
	}
}
