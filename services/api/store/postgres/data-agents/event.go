package dataagents

import (
	"context"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	pg "github.com/consensys/orchestrate/pkg/toolkit/database/postgres"
	"github.com/consensys/orchestrate/services/api/store"
	"github.com/consensys/orchestrate/services/api/store/models"
)

const eventDAComponent = "data-agents.event"

// PGAccount is an Account data agent for PostgreSQL
type PGEvent struct {
	db     pg.DB
	logger *log.Logger
}

// NewPGAccount creates a new PGAccount
func NewPGEvent(db pg.DB) store.EventAgent {
	return &PGEvent{db: db, logger: log.NewLogger().SetComponent(eventDAComponent)}
}

func (agent *PGEvent) InsertMultiple(ctx context.Context, events []*models.EventModel) error {
	query := agent.db.ModelContext(ctx, &events).OnConflict("DO NOTHING")
	err := pg.InsertQuery(ctx, query)
	if err != nil {
		return errors.FromError(err).ExtendComponent(eventDAComponent)
	}

	return nil
}

func (agent *PGEvent) FindOneByAccountAndSigHash(ctx context.Context, chainID, address, sighash string, indexedInputCount uint32) (*models.EventModel, error) {
	event := &models.EventModel{}
	query := agent.db.ModelContext(ctx, event).
		Column("event_model.abi").
		Join("JOIN codehashes AS c ON c.codehash = event_model.codehash").
		Where("c.chain_id = ?", chainID).
		Where("c.address = ?", address).
		Where("event_model.sig_hash = ?", sighash).
		Where("event_model.indexed_input_count = ?", indexedInputCount)

	err := pg.SelectOne(ctx, query)
	if err != nil {
		if !errors.IsNotFoundError(err) {
			agent.logger.WithContext(ctx).WithError(err).Error("failed to find event")
		}
		return nil, errors.FromError(err).ExtendComponent(eventDAComponent)
	}

	return event, nil
}
func (agent *PGEvent) FindDefaultBySigHash(ctx context.Context, sighash string, indexedInputCount uint32) ([]*models.EventModel, error) {
	var defaultEvents []*models.EventModel
	query := agent.db.ModelContext(ctx, &defaultEvents).
		ColumnExpr("DISTINCT abi").
		Where("sig_hash = ?", sighash).
		Where("indexed_input_count = ?", indexedInputCount).
		Order("abi DESC")

	err := pg.Select(ctx, query)
	if err != nil {
		if !errors.IsNotFoundError(err) {
			agent.logger.WithContext(ctx).WithError(err).Error("failed to find default event")
		}
		return nil, errors.FromError(err).ExtendComponent(eventDAComponent)
	}

	return defaultEvents, nil
}
