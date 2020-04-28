package store

import (
	"context"
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/interfaces"

	postgresStore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/postgres"
)

const storeServiceName = "transaction-scheduler.store"

func Build(ctx context.Context, cfg *Config, pgmngr postgres.Manager) (interfaces.DB, error) {
	switch cfg.Type {
	case postgresType:
		cfg.Postgres.PG.ApplicationName = storeServiceName
		opts, err := cfg.Postgres.PG.PGOptions()
		if err != nil {
			return nil, err
		}
		return postgresStore.NewPGDB(pgmngr.Connect(ctx, opts)), nil
	default:
		return nil, fmt.Errorf("invalid transaction scheduler store type %q", cfg.Type)
	}
}
