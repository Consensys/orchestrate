package store

import (
	"context"
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
)

const storeServiceName = "transaction-scheduler.store"

func Build(ctx context.Context, cfg *Config, pgmngr postgres.Manager) (*DataAgents, error) {
	switch cfg.Type {
	case postgresType:
		cfg.Postgres.PG.ApplicationName = storeServiceName
		db := pgmngr.Connect(ctx, cfg.Postgres.PG)

		return &DataAgents{
			TransactionRequest: db,
		}, nil
	default:
		return &DataAgents{}, fmt.Errorf("invalid transaction scheduler store type %q", cfg.Type)
	}
}
