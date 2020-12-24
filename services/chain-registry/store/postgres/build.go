package postgres

import (
	"github.com/go-pg/pg/v9"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store"
	pgda "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store/postgres/data-agents"
)

func Build(db *pg.DB) store.DataAgents {
	return store.DataAgents{
		Chain: pgda.NewPGChainAgent(db),
	}
}
