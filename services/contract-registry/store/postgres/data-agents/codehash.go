package dataagents

import (
	"context"

	"github.com/go-pg/pg/v9"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/models"
)

// PGCodeHash is a codehash data agent
type PGCodeHash struct {
	db *pg.DB
}

// NewPGCodeHash creates a new PGCodeHash
func NewPGCodeHash(db *pg.DB) *PGCodeHash {
	return &PGCodeHash{db: db}
}

// Insert Inserts a new codehash in DB
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
		return errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	return nil
}
