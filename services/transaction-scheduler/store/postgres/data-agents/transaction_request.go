package dataagents

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"

	"github.com/go-pg/pg/v9"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

const txRequestDAComponent = "data-agents.transaction-request"

// PGTransactionRequest is a transaction request data agent for PostgreSQL
type PGTransactionRequest struct {
	db *pg.DB
}

// NewPGTransactionRequest creates a new PGTransactionRequest
func NewPGTransactionRequest(db *pg.DB) *PGTransactionRequest {
	return &PGTransactionRequest{db: db}
}

// Insert Inserts a new transaction request in DB
func (agent *PGTransactionRequest) SelectOrInsert(ctx context.Context, txRequest *models.TransactionRequest) error {
	_, err := agent.db.ModelContext(ctx, txRequest).
		Where("idempotency_key = ?idempotency_key").
		OnConflict("ON CONSTRAINT requests_idempotency_key_key DO NOTHING").
		SelectOrInsert()
	if err != nil {
		pgErr, ok := err.(pg.Error)
		if ok && pgErr.IntegrityViolation() || errors.IsAlreadyExistsError(err) {
			errMessage := "transaction request already exists"
			log.WithContext(ctx).WithError(err).Error(errMessage)
			return errors.AlreadyExistsError(errMessage).ExtendComponent(txRequestDAComponent)
		}

		errMessage := "error executing selectOrInsert"
		log.WithContext(ctx).WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(txRequestDAComponent)
	}

	return nil
}

func (agent *PGTransactionRequest) FindOneByIdempotencyKey(ctx context.Context, idempotencyKey string) (*models.TransactionRequest, error) {
	txRequest := &models.TransactionRequest{}
	err := agent.db.ModelContext(ctx, txRequest).
		Where("idempotency_key = ?", idempotencyKey).
		First()

	if err != nil && err == pg.ErrNoRows {
		errMessage := "could not load transaction request with idempotency_key: %s"
		log.WithError(err).Debugf(errMessage, idempotencyKey)
		return nil, errors.NotFoundError(errMessage, idempotencyKey).ExtendComponent(txRequestDAComponent)
	} else if err != nil {
		errMessage := "Failed to get transaction request from DB"
		log.WithError(err).Error(errMessage)
		return nil, errors.PostgresConnectionError(errMessage).ExtendComponent(txRequestDAComponent)
	}

	return txRequest, nil
}
