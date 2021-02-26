package dataagents

import (
	"context"
	"time"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
	"github.com/ConsenSys/orchestrate/services/api/store"

	pg "github.com/ConsenSys/orchestrate/pkg/toolkit/database/postgres"
	"github.com/ConsenSys/orchestrate/services/api/store/models"
	"github.com/gofrs/uuid"
)

const txDAComponent = "transaction.log"

// PGLog is a log data agent for PostgreSQL
type PGTransaction struct {
	db     pg.DB
	logger *log.Logger
}

// NewPGLog creates a new PGLog
func NewPGTransaction(db pg.DB) store.TransactionAgent {
	return &PGTransaction{db: db, logger: log.NewLogger().SetComponent(txDAComponent)}
}

// Insert Inserts a new log in DB
func (agent *PGTransaction) Insert(ctx context.Context, txModel *models.Transaction) error {
	if txModel.UUID == "" {
		txModel.UUID = uuid.Must(uuid.NewV4()).String()
	}

	err := pg.Insert(ctx, agent.db, txModel)
	if err != nil {
		agent.logger.WithContext(ctx).WithError(err).Error("failed to insert transaction")
		return errors.FromError(err).ExtendComponent(txDAComponent)
	}

	return nil
}

// Insert Inserts a new log in DB
func (agent *PGTransaction) Update(ctx context.Context, txModel *models.Transaction) error {
	if txModel.ID == 0 {
		err := errors.InvalidArgError("cannot update transaction with missing ID")
		agent.logger.WithContext(ctx).WithError(err).Error("failed to insert transaction")
		return err
	}

	txModel.UpdatedAt = time.Now().UTC()
	err := pg.Update(ctx, agent.db, txModel)
	if err != nil {
		agent.logger.WithContext(ctx).WithError(err).Error("failed to update transaction")
		return errors.FromError(err).ExtendComponent(txDAComponent)
	}

	return nil
}
