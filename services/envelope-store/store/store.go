package store

import (
	"context"
	"fmt"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
	memorystore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store/memory"
	pgstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store/postgres"
)

const StoreName = "envelopes"

//go:generate mockgen -source=store.go -destination=mock/mock.go -package=mock

type Builder interface {
	Build(ctx context.Context, conf *Config) (svc.EnvelopeStoreServer, error)
}

func NewBuilder(pgmngr postgres.Manager) Builder {
	return &builder{
		postgres: pgstore.NewBuilder(pgmngr),
		memory:   memorystore.NewBuilder(),
	}
}

type builder struct {
	postgres *pgstore.Builder
	memory   *memorystore.Builder
}

func (b *builder) Build(ctx context.Context, conf *Config) (svc.EnvelopeStoreServer, error) {
	logCtx := log.With(ctx, log.Str("store", StoreName))
	switch conf.Type {
	case postgresType:
		conf.Postgres.PG.ApplicationName = StoreName
		return b.postgres.Build(logCtx, conf.Postgres)
	case inMemoryType:
		return b.memory.Build(logCtx)
	default:
		return nil, fmt.Errorf("invalid envelope registry store type %q", conf.Type)
	}
}
