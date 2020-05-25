package dataagents

import (
	"context"

	log "github.com/sirupsen/logrus"
	pg "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
)

const txRequestDAComponent = "data-agents.transaction-request"

// PGTransactionRequest is a transaction request data agent for PostgreSQL
type PGTransactionRequest struct {
	db pg.DB
}

// NewPGTransactionRequest creates a new PGTransactionRequest
func NewPGTransactionRequest(db pg.DB) *PGTransactionRequest {
	return &PGTransactionRequest{db: db}
}

// Insert Inserts a new transaction request in DB
func (agent *PGTransactionRequest) SelectOrInsert(ctx context.Context, txRequest *models.TransactionRequest) error {
	_, err := agent.db.ModelContext(ctx, txRequest).
		Where("idempotency_key = ?idempotency_key").
		OnConflict("ON CONSTRAINT transaction_requests_idempotency_key_key DO NOTHING").
		Relation("Schedules").
		SelectOrInsert()

	if err != nil {
		errMessage := "error executing selectOrInsert"
		log.WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(txRequestDAComponent)
	}

	return nil
}

func (agent *PGTransactionRequest) FindOneByIdempotencyKey(ctx context.Context, idempotencyKey string) (*models.TransactionRequest, error) {
	txRequest := &models.TransactionRequest{}
	query := agent.db.ModelContext(ctx, txRequest).
		Where("idempotency_key = ?", idempotencyKey).
		Relation("Schedules")

	err := pg.SelectOne(ctx, query)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(jobDAComponent)
	}

	return txRequest, nil
}
