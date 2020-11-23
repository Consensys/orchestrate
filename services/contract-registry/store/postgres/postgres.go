package postgres

import (
	"github.com/go-pg/pg/v9"
	dataagents "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/contract-registry/store/postgres/data-agents"
)

func Build(db *pg.DB) (
	*dataagents.PGContract,
	*dataagents.PGRepository,
	*dataagents.PGTag,
	*dataagents.PGArtifact,
	*dataagents.PGMethod,
	*dataagents.PGEvent,
	*dataagents.PGCodeHash,
) {
	repositories := dataagents.NewPGRepository(db)
	artifacts := dataagents.NewPGArtifact(db)
	codeHashes := dataagents.NewPGCodeHash(db)
	events := dataagents.NewPGEvent(db)
	methods := dataagents.NewPGMethod(db)
	tags := dataagents.NewPGTag(db)
	contracts := dataagents.NewPGContract(db, repositories, artifacts, tags, methods, events)

	return contracts, repositories, tags, artifacts, methods, events, codeHashes
}
