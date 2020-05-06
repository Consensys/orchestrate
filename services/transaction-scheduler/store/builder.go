package store

import (
	"context"
	"fmt"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	pgstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/postgres"
)

const storeServiceName = "transaction-scheduler.store"

type Builder interface {
	Build(context.Context, *Config) (*DataAgents, error)
}

func NewBuilder(pgmngr postgres.Manager) Builder {
	return &builder{
		postgres: pgstore.NewBuilder(pgmngr),
	}
}

type builder struct {
	postgres *pgstore.Builder
}

func (b *builder) Build(ctx context.Context, cfg *Config) (*DataAgents, error) {
	logCtx := log.With(ctx, log.Str("store", storeServiceName))
	switch cfg.Type {
	case postgresType:
		cfg.Postgres.PG.ApplicationName = storeServiceName
		schedules, jobs, logs, txRequests, err := b.postgres.Build(logCtx, cfg.Postgres)
		return &DataAgents{
			ScheduleAgent:      schedules,
			JobAgent:           jobs,
			LogAgent:           logs,
			TransactionRequest: txRequests,
		}, err
	default:
		return &DataAgents{}, fmt.Errorf("invalid transaction scheduler store type %q", cfg.Type)
	}
}
