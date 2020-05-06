package postgres

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	pgda "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/postgres/data-agents"
)

type Builder struct {
	postgres postgres.Manager
}

func NewBuilder(mngr postgres.Manager) *Builder {
	return &Builder{postgres: mngr}
}

func (b *Builder) Build(ctx context.Context, cfg *Config) (*pgda.PGSchedule, *pgda.PGJob, *pgda.PGLog, *pgda.PGTransactionRequest, error) {
	db := b.postgres.Connect(ctx, cfg.PG)
	return pgda.NewPGSchedule(db), pgda.NewPGJob(db), pgda.NewPGLog(db), pgda.NewPGTransactionRequest(db), nil
}
