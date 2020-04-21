package store

import (
	"context"
	"fmt"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	pgstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/postgres"
)

const storeServiceName = "chains"

type Builder interface {
	Build(context.Context, *Config) (DataAgents, error)
}

func NewBuilder(pgmngr postgres.Manager) Builder {
	return &builder{
		postgres: pgstore.NewBuilder(pgmngr),
	}
}

type builder struct {
	postgres *pgstore.Builder
}

func (b *builder) Build(ctx context.Context, conf *Config) (DataAgents, error) {
	logCtx := log.With(ctx, log.Str("store", storeServiceName))
	switch conf.Type {
	case postgresType:
		conf.Postgres.PG.ApplicationName = storeServiceName
		chainAgent, fauceAgent, privateTxAgent, err := b.postgres.Build(logCtx, conf.Postgres)
		return DataAgents{
			Chain:     chainAgent,
			Faucet:    fauceAgent,
			PrivateTx: privateTxAgent,
		}, err
	default:
		return DataAgents{}, fmt.Errorf("invalid chain registry store type %q", conf.Type)
	}
}
