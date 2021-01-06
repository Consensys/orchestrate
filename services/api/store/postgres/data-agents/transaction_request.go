package dataagents

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"

	gopg "github.com/go-pg/pg/v9"
	log "github.com/sirupsen/logrus"
	pg "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models"
)

const txRequestDAComponent = "data-agents.transaction-request"

// PGTransactionRequest is a transaction request data agent for PostgreSQL
type PGTransactionRequest struct {
	db pg.DB
}

// NewPGTransactionRequest creates a new PGTransactionRequest
func NewPGTransactionRequest(db pg.DB) store.TransactionRequestAgent {
	return &PGTransactionRequest{db: db}
}

// Insert Inserts a new transaction request in DB
func (agent *PGTransactionRequest) Insert(ctx context.Context, txRequest *models.TransactionRequest) error {
	if txRequest.Schedule != nil && txRequest.ScheduleID == nil {
		txRequest.ScheduleID = &txRequest.Schedule.ID
	}

	err := pg.Insert(ctx, agent.db, txRequest)
	if err != nil {
		errMessage := "error executing selectOrInsert"
		log.WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(txRequestDAComponent)
	}

	return nil
}

func (agent *PGTransactionRequest) FindOneByIdempotencyKey(ctx context.Context, idempotencyKey, tenantID string) (*models.TransactionRequest, error) {
	txRequest := &models.TransactionRequest{}
	query := agent.db.ModelContext(ctx, txRequest).
		Where("idempotency_key = ?", idempotencyKey).
		Relation("Schedule")

	query = query.Where("schedule.tenant_id = ?", tenantID)

	err := pg.SelectOne(ctx, query)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(txRequestDAComponent)
	}

	return txRequest, nil
}

func (agent *PGTransactionRequest) FindOneByUUID(ctx context.Context, scheduleUUID string, tenants []string) (*models.TransactionRequest, error) {
	txRequest := &models.TransactionRequest{}
	query := agent.db.ModelContext(ctx, txRequest).
		Where("schedule.uuid = ?", scheduleUUID).
		Relation("Schedule")

	query = pg.WhereAllowedTenants(query, "schedule.tenant_id", tenants)

	err := pg.SelectOne(ctx, query)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(txRequestDAComponent)
	}

	return txRequest, nil
}

func (agent *PGTransactionRequest) Search(ctx context.Context, filters *entities.TransactionRequestFilters, tenants []string) ([]*models.TransactionRequest, error) {
	var txRequests []*models.TransactionRequest

	query := agent.db.ModelContext(ctx, &txRequests).Relation("Schedule")

	if len(filters.IdempotencyKeys) > 0 {
		query = query.Where("transaction_request.idempotency_key in (?)", gopg.In(filters.IdempotencyKeys))
	}

	query = pg.WhereAllowedTenants(query, "schedule.tenant_id", tenants)

	err := pg.Select(ctx, query)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(txRequestDAComponent)
	}

	return txRequests, nil
}
