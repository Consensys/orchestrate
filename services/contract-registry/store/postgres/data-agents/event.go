package dataagents

import (
	"context"

	"github.com/go-pg/pg/v9/orm"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"

	"github.com/go-pg/pg/v9"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/common"
)

// PGEvent is an event data agent
type PGEvent struct {
	db *pg.DB
}

// NewPGEvent creates a new PGEvent
func NewPGEvent(db *pg.DB) *PGEvent {
	return &PGEvent{db: db}
}

// InsertMultiple Inserts multiple new events in DB
func (agent *PGEvent) InsertMultiple(ctx context.Context, events *[]*models.EventModel) error {
	tx := postgres.TxFromContext(ctx)
	if tx != nil {
		return agent.insertMultiple(tx.ModelContext(ctx, events))
	}

	return agent.insertMultiple(agent.db.ModelContext(ctx, events))

}

func (agent *PGEvent) insertMultiple(query *orm.Query) error {
	_, err := query.
		OnConflict("DO NOTHING").
		Insert()
	if err != nil {
		errMessage := "could not create events"
		log.WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	return nil
}

// FindOneByAccountAndSigHash Finds an event by account and sighash
func (agent *PGEvent) FindOneByAccountAndSigHash(ctx context.Context, account *common.AccountInstance, sighash string, indexedInputCount uint32) (*models.EventModel, error) {
	event := &models.EventModel{}
	err := agent.db.ModelContext(ctx, event).
		Column("event_model.abi").
		Join("JOIN codehashes AS c ON c.codehash = event_model.codehash").
		Where("c.chain_id = ?", account.GetChainId()).
		Where("c.address = ?", account.GetAccount()).
		Where("event_model.sig_hash = ?", sighash).
		Where("event_model.indexed_input_count = ?", indexedInputCount).
		First()

	if err != nil && err == pg.ErrNoRows {
		errMessage := "could not load event with chainId: %s, account: %s, sighash %s and indexedInputCount %v"
		log.WithError(err).Debugf(errMessage, account.GetChainId(), account.GetAccount(), sighash, indexedInputCount)
		return nil, errors.NotFoundError(errMessage, account.GetChainId(), account.GetAccount(), sighash, indexedInputCount).ExtendComponent(component)
	} else if err != nil {
		errMessage := "Failed to get event from DB"
		log.WithError(err).Error(errMessage)
		return nil, errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	return event, nil
}

// FindDefaultBySigHash Find events by sighash
func (agent *PGEvent) FindDefaultBySigHash(ctx context.Context, sighash string, indexedInputCount uint32) ([]*models.EventModel, error) {
	var defaultEvents []*models.EventModel
	err := agent.db.ModelContext(ctx, &defaultEvents).
		ColumnExpr("DISTINCT abi").
		Where("sig_hash = ?", sighash).
		Where("indexed_input_count = ?", indexedInputCount).
		Order("abi DESC").
		Select()

	if err != nil {
		errMessage := "Failed to get default events from DB"
		log.WithError(err).Error(errMessage)
		return nil, errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	if len(defaultEvents) == 0 {
		errMessage := "could not load default events with sighash: %s and indexedInputCount %v"
		log.WithError(err).Debugf(errMessage, sighash, indexedInputCount)
		return nil, errors.NotFoundError(errMessage, sighash, indexedInputCount).ExtendComponent(component)
	}

	return defaultEvents, nil
}
