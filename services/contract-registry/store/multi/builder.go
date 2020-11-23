package multi

import (
	"context"

	"github.com/go-pg/pg/v9"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres"
)

const storeServiceName = "contract-registry.store"

func Build(ctx context.Context, cfg *Config, pgmngr postgres.Manager) (*pg.DB, error) {
	cfg.Postgres.PG.ApplicationName = storeServiceName
	opts, err := cfg.Postgres.PG.PGOptions()
	if err != nil {
		return nil, err
	}
	return pgmngr.Connect(ctx, opts), nil
}
