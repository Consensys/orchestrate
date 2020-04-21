package dataagents

import (
	"context"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

const privateTxComponentName = "chain-registry.store.privateTx.pg"

// PGPrivateTxAgent is an artifact data agent
type PGPrivateTxAgent struct {
	db *pg.DB
}

// NewPGArtifact creates a new PGPrivateTxAgent
func NewPGPrivateTxManager(db *pg.DB) *PGPrivateTxAgent {
	return &PGPrivateTxAgent{db: db}
}

// SelectOrInsert Inserts a new artifact in DB
func (agent *PGPrivateTxAgent) InsertMultiple(ctx context.Context, privateTxManager *[]*models.PrivateTxManagerModel) error {
	if tx := postgres.TxFromContext(ctx); tx != nil {
		return agent.insertMultiple(tx.ModelContext(ctx, privateTxManager))
	}

	return agent.insertMultiple(agent.db.ModelContext(ctx, privateTxManager))
}

func (agent *PGPrivateTxAgent) insertMultiple(query *orm.Query) error {
	_, err := query.
		Insert()
	if err != nil {
		errMessage := "could not create private tx manager"
		log.WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(privateTxComponentName)
	}

	return nil
}
