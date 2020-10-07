package multi

import (
	"context"
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/store"
	storePG "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/store/postgres"
)

const storeServiceName = "identity-manager.store"

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
		return nil, fmt.Errorf("invalid identity manager store type %q", cfg.Type)
	}
}
