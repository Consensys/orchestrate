package multi

import (
	"context"
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/store"
	storePG "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/store/postgres"
)

const storeServiceName = "transaction-scheduler.store"

func Build(ctx context.Context, cfg *Config, pgmngr postgres.Manager) (store.DB, error) {
	switch cfg.Type {
	case postgresType:
		cfg.Postgres.PG.ApplicationName = storeServiceName
		opts, err := cfg.Postgres.PG.PGOptions()
		if err != nil {
			return nil, err
		}
		return storePG.NewPGDB(pgmngr.Connect(ctx, opts)), nil
	default:
		return nil, fmt.Errorf("invalid transaction scheduler store type %q", cfg.Type)
	}
}
