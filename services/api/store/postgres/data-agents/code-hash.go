package dataagents

import (
	"context"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	pkgpg "github.com/consensys/orchestrate/pkg/toolkit/database/postgres"
	"github.com/consensys/orchestrate/services/api/store"
	"github.com/consensys/orchestrate/services/api/store/models"
)

const codeHashDAComponent = "data-agents.code_hash"

type PGCodeHash struct {
	db     pkgpg.DB
	logger *log.Logger
}

func NewPGCodeHash(db pkgpg.DB) store.CodeHashAgent {
	return &PGCodeHash{db: db, logger: log.NewLogger().SetComponent(codeHashDAComponent)}
}

func (agent *PGCodeHash) Insert(ctx context.Context, codehash *models.CodehashModel) error {
	// If uniqueness constraint is broken then it updates the former value
	_, err := agent.db.ModelContext(ctx, codehash).
		OnConflict("ON CONSTRAINT codehashes_chain_id_address_key DO UPDATE").
		Set("chain_id = ?chain_id").
		Set("address = ?address").
		Set("codehash = ?codehash").
		Returning("*").
		Insert()

	if err != nil {
		errMessage := "could not insert codehash"
		agent.logger.WithContext(ctx).WithError(err).Error(errMessage)
		return errors.FromError(err).ExtendComponent(codeHashDAComponent)
	}

	return nil
}
