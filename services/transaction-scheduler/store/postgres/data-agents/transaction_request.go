package dataagents

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"

	"github.com/go-pg/pg/v9/orm"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"

	"github.com/go-pg/pg/v9"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

const txRequestDAComponent = "data-agents.transaction-request"

// PGRepository is a repository data agent
type PGTransactionRequest struct {
	db *pg.DB
}

// NewPGTransactionRequest creates a new PGTransactionRequest
func NewPGTransactionRequest(db *pg.DB) *PGTransactionRequest {
	return &PGTransactionRequest{db: db}
}

// TODO: Fix the way we pass the tx from postgres to children data agents
// Insert Inserts a new transaction request in DB
func (agent *PGTransactionRequest) Insert(ctx context.Context, txRequest *models.TransactionRequest) error {
	tx := postgres.TxFromContext(ctx)
	if tx != nil {
		return agent.insert(ctx, tx.ModelContext(ctx, txRequest))
	}

	return agent.insert(ctx, agent.db.ModelContext(ctx, txRequest))
}

func (agent *PGTransactionRequest) insert(ctx context.Context, query *orm.Query) error {
	_, err := query.Insert()
	if err != nil {
		pgErr, ok := err.(pg.Error)
		if ok && pgErr.IntegrityViolation() || errors.IsAlreadyExistsError(err) {
			errMessage := "transaction request already exists"
			log.WithContext(ctx).WithError(err).Error(errMessage)
			return errors.AlreadyExistsError(errMessage).ExtendComponent(txRequestDAComponent)
		}

		errMessage := "error executing insert"
		log.WithContext(ctx).WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(txRequestDAComponent)
	}

	return nil
}
