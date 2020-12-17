package multi

import (
	"context"
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
	storePG "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/postgres"
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
