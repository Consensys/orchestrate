package dataagents

import (
	"context"

	log "github.com/sirupsen/logrus"
	pkgpg "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models"
)

const codeHashDAComponent = "data-agents.code_hash"

type PGCodeHash struct {
	db pkgpg.DB
}

func NewPGCodeHash(db pkgpg.DB) store.CodeHashAgent {
	return &PGCodeHash{db: db}
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
		errMessage := "could not create codehash"
		log.WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(codeHashDAComponent)
	}

	return nil
}
