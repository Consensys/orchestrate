package dataagents

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"

	uuid "github.com/satori/go.uuid"
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
	tx, err := agent.db.Begin()
	if err != nil {
		errMessage := "failed to create DB transaction"
		log.WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(jobDAComponent)
	}

	// Insert schedule
	schedule := txRequest.Schedule
	schedule.UUID = uuid.NewV4().String()
	err = postgres.Insert(ctx, tx, schedule)
	if err != nil {
		return errors.FromError(err).ExtendComponent(txRequestDAComponent)
	}

	// SelectOrInsert TransactionRequest
	txRequest.ScheduleID = schedule.ID
	_, err = tx.ModelContext(ctx, txRequest).
		Where("idempotency_key = ?idempotency_key").
		OnConflict("ON CONSTRAINT requests_idempotency_key_key DO NOTHING").
		SelectOrInsert()
	if err != nil {
		errMessage := "error executing selectOrInsert"
		log.WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(txRequestDAComponent)
	}

	// Insert job
	for _, job := range txRequest.Schedule.Jobs {
		transaction := job.Transaction
		transaction.UUID = uuid.NewV4().String()
		err = postgres.Insert(ctx, tx, transaction)
		if err != nil {
			return errors.FromError(err).ExtendComponent(txRequestDAComponent)
		}

		job.TransactionID = transaction.ID
		job.ScheduleID = schedule.ID
		job.UUID = uuid.NewV4().String()
		err = postgres.Insert(ctx, tx, job)
		if err != nil {
			return errors.FromError(err).ExtendComponent(txRequestDAComponent)
		}

		for _, logModel := range job.Logs {
			logModel.UUID = uuid.NewV4().String()
			logModel.JobID = job.ID
			err = postgres.Insert(ctx, tx, logModel)
			if err != nil {
				return errors.FromError(err).ExtendComponent(txRequestDAComponent)
			}
		}
	}

	return tx.Commit()
}

func (agent *PGTransactionRequest) FindOneByIdempotencyKey(ctx context.Context, idempotencyKey string) (*models.TransactionRequest, error) {
	txRequest := &models.TransactionRequest{}
	err := agent.db.ModelContext(ctx, txRequest).
		Where("idempotency_key = ?", idempotencyKey).
		Relation("Schedule").
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
