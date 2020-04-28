package postgres

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	pgda "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store/postgres/data-agents"
)

type Builder struct {
	postgres postgres.Manager
}

func NewBuilder(mngr postgres.Manager) *Builder {
	return &Builder{postgres: mngr}
}

func (b *Builder) Build(ctx context.Context, cfg *Config) (*pgda.PGEnvelopeAgent, error) {
	opts, err := cfg.PG.PGOptions()
	if err != nil {
		return nil, err
	}
	db := b.postgres.Connect(ctx, opts)

	return pgda.NewPGEnvelope(db), nil
}
