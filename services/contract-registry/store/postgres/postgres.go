package postgres

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	dataagents "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/postgres/data-agents"
)

type Builder struct {
	postgres postgres.Manager
}

func NewBuilder(mngr postgres.Manager) *Builder {
	return &Builder{postgres: mngr}
}

func (b *Builder) Build(ctx context.Context, cfg *Config) (
	*dataagents.PGContract,
	*dataagents.PGRepository,
	*dataagents.PGTag,
	*dataagents.PGArtifact,
	*dataagents.PGMethod,
	*dataagents.PGEvent,
	*dataagents.PGCodeHash,
	error,
) {
	db := b.postgres.Connect(ctx, cfg.PG)

	repositories := dataagents.NewPGRepository(db)
	artifacts := dataagents.NewPGArtifact(db)
	codeHashes := dataagents.NewPGCodeHash(db)
	events := dataagents.NewPGEvent(db)
	methods := dataagents.NewPGMethod(db)
	tags := dataagents.NewPGTag(db)
	contracts := dataagents.NewPGContract(db, repositories, artifacts, tags, methods, events)

	return contracts, repositories, tags, artifacts, methods, events, codeHashes, nil
}
