package store

import (
	"context"
	"fmt"

	dataagents "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/postgres/data-agents"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
)

const storeServiceName = "transaction-scheduler.store"

func Build(ctx context.Context, cfg *Config, pgmngr postgres.Manager) (*DataAgents, error) {
	switch cfg.Type {
	case postgresType:
		cfg.Postgres.PG.ApplicationName = storeServiceName
		db := pgmngr.Connect(ctx, cfg.Postgres.PG)

		return &DataAgents{
			TransactionRequest: dataagents.NewPGTransactionRequest(db),
			ScheduleAgent:      dataagents.NewPGSchedule(db),
			JobAgent:           dataagents.NewPGJob(db),
			LogAgent:           dataagents.NewPGLog(db),
		}, nil
	default:
		return &DataAgents{}, fmt.Errorf("invalid transaction scheduler store type %q", cfg.Type)
	}
}
