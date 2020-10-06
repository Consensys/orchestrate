package multi

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v9"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
)

const storeServiceName = "chain-registry.store"

func Build(ctx context.Context, cfg *Config, pgmngr postgres.Manager) (*pg.DB, error) {
	switch cfg.Type {
	case postgresType:
		cfg.Postgres.PG.ApplicationName = storeServiceName
		opts, err := cfg.Postgres.PG.PGOptions()
		if err != nil {
			return nil, err
		}
		return pgmngr.Connect(ctx, opts), nil
	default:
		return nil, fmt.Errorf("invalid transaction scheduler store type %q", cfg.Type)
	}
}
