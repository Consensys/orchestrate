package multi

import (
	"context"
	"fmt"

	"github.com/ConsenSys/orchestrate/pkg/toolkit/database/postgres"
	"github.com/ConsenSys/orchestrate/services/api/store"
	storePG "github.com/ConsenSys/orchestrate/services/api/store/postgres"
)

const storeServiceName = "api.store"

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
		return nil, fmt.Errorf("invalid API store type %q", cfg.Type)
	}
}
