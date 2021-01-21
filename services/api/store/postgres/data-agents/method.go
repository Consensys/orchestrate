package dataagents

import (
	"context"

	pg "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models"
)

const methodDAComponent = "data-agents.method"

// PGAccount is an Account data agent for PostgreSQL
type PGMethod struct {
	db     pg.DB
	logger *log.Logger
}

// NewPGAccount creates a new PGAccount
func NewPGMethod(db pg.DB) store.MethodAgent {
	return &PGMethod{db: db, logger: log.NewLogger().SetComponent(methodDAComponent)}
}

func (agent *PGMethod) InsertMultiple(ctx context.Context, methods []*models.MethodModel) error {
	query := agent.db.ModelContext(ctx, &methods).
		OnConflict("DO NOTHING")

	err := pg.InsertQuery(ctx, query)
	if err != nil {
		agent.logger.WithContext(ctx).WithError(err).Error("failed to insert multiple contract methods")
		return errors.FromError(err).ExtendComponent(methodDAComponent)
	}

	return nil
}
func (agent *PGMethod) FindOneByAccountAndSelector(ctx context.Context, chainID, address string, selector []byte) (*models.MethodModel, error) {
	method := &models.MethodModel{}
	query := agent.db.ModelContext(ctx, method).
		Column("method_model.abi").
		Join("JOIN codehashes AS c ON c.codehash = method_model.codehash").
		Where("c.chain_id = ?", chainID).
		Where("c.address = ?", address).
		Where("method_model.selector = ?", selector)

	err := pg.SelectOne(ctx, query)
	if err != nil {
		if !errors.IsNotFoundError(err) {
			agent.logger.WithContext(ctx).WithError(err).Error("failed to find contract method by account")
		}
		return nil, errors.FromError(err).ExtendComponent(methodDAComponent)
	}

	return method, nil
}
func (agent *PGMethod) FindDefaultBySelector(ctx context.Context, selector []byte) ([]*models.MethodModel, error) {
	var defaultMethods []*models.MethodModel
	query := agent.db.ModelContext(ctx, &defaultMethods).
		ColumnExpr("DISTINCT abi").
		Where("selector = ?", selector)

	err := pg.Select(ctx, query)
	if err != nil {
		if !errors.IsNotFoundError(err) {
			agent.logger.WithContext(ctx).WithError(err).Error("failed to find default contract method")
		}
		return nil, errors.FromError(err).ExtendComponent(methodDAComponent)
	}

	return defaultMethods, nil
}
