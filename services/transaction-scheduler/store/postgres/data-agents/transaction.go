package dataagents

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"

	"github.com/gofrs/uuid"
	pg "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
)

const txDAComponent = "transaction.log"

// PGLog is a log data agent for PostgreSQL
type PGTransaction struct {
	db pg.DB
}

// NewPGLog creates a new PGLog
func NewPGTransaction(db pg.DB) *PGTransaction {
	return &PGTransaction{db: db}
}

// Insert Inserts a new log in DB
func (agent *PGTransaction) Insert(ctx context.Context, txModel *models.Transaction) error {
	if txModel.UUID == "" {
		txModel.UUID = uuid.Must(uuid.NewV4()).String()
	}

	err := pg.Insert(ctx, agent.db, txModel)
	if err != nil {
		return errors.FromError(err).ExtendComponent(txDAComponent)
	}

	return nil
}

// Insert Inserts a new log in DB
func (agent *PGTransaction) Update(ctx context.Context, txModel *models.Transaction) error {
	if txModel.ID == 0 {
		errMsg := "cannot update transaction with missing ID"
		log.WithContext(ctx).Error(errMsg)
		return errors.InvalidArgError(errMsg)
	}

	err := pg.Update(ctx, agent.db, txModel)
	if err != nil {
		return errors.FromError(err).ExtendComponent(txDAComponent)
	}

	return nil
}
