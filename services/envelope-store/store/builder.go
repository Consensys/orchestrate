package store

import (
	"context"
	"fmt"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	pgstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store/postgres"
)

const storeServiceName = "envelopes"

type Builder interface {
	Build(ctx context.Context, conf *Config) (DataAgents, error)
}

type builder struct {
	postgres *pgstore.Builder
}

func NewBuilder(pgmngr postgres.Manager) Builder {
	return &builder{
		postgres: pgstore.NewBuilder(pgmngr),
	}
}

func (b *builder) Build(ctx context.Context, cfg *Config) (DataAgents, error) {
	logCtx := log.With(ctx, log.Str("store", storeServiceName))
	switch cfg.Type {
	case postgresType:
		cfg.Postgres.PG.ApplicationName = storeServiceName
		pgEnvelopeAgent, err := b.postgres.Build(logCtx, cfg.Postgres)
		return DataAgents{
			Envelope: pgEnvelopeAgent,
		}, err
	default:
		return DataAgents{}, fmt.Errorf("invalid envelope registry store type %q", cfg.Type)
	}
}
