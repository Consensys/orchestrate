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

// PGMethod is a method data agent
type PGMethod struct {
	db *pg.DB
}

// NewPGMethod creates a new PGMethod
func NewPGMethod(db *pg.DB) *PGMethod {
	return &PGMethod{db: db}
}

// InsertMultiple Inserts multiple new methods in DB
func (agent *PGMethod) InsertMultiple(ctx context.Context, methods *[]*models.MethodModel) error {
	tx := postgres.TxFromContext(ctx)
	if tx != nil {
		return agent.insertMultiple(tx.ModelContext(ctx, methods))
	}

	return agent.insertMultiple(agent.db.ModelContext(ctx, methods))
}

func (agent *PGMethod) insertMultiple(query *orm.Query) error {
	_, err := query.
		OnConflict("DO NOTHING").
		Insert()
	if err != nil {
		errMessage := "could not create methods"
		log.WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	return nil
}

// Finds a method by account and selector
func (agent *PGMethod) FindOneByAccountAndSelector(ctx context.Context, account *common.AccountInstance, selector []byte) (*models.MethodModel, error) {
	method := &models.MethodModel{}
	err := agent.db.ModelContext(ctx, method).
		Column("method_model.abi").
		Join("JOIN codehashes AS c ON c.codehash = method_model.codehash").
		Where("c.chain_id = ?", account.GetChainId()).
		Where("c.address = ?", account.GetAccount()).
		Where("method_model.selector = ?", selector).
		First()

	if err != nil && err == pg.ErrNoRows {
		errMessage := "could not load method with chainId: %s, account: %s and selector %v"
		log.WithError(err).Debugf(errMessage, account.GetChainId(), account.GetAccount(), selector)
		return nil, errors.NotFoundError(errMessage, account.GetChainId(), account.GetAccount(), selector).ExtendComponent(component)
	} else if err != nil {
		errMessage := "Failed to get method from DB"
		log.WithError(err).Error(errMessage)
		return nil, errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	return method, nil
}

// FindDefaultBySelector Finds methods by selector
func (agent *PGMethod) FindDefaultBySelector(ctx context.Context, selector []byte) ([]*models.MethodModel, error) {
	var defaultMethods []*models.MethodModel
	err := agent.db.ModelContext(ctx, &defaultMethods).
		ColumnExpr("DISTINCT abi").
		Where("selector = ?", selector).
		Select()

	if err != nil {
		errMessage := "Failed to get default methods from DB"
		log.WithError(err).Error(errMessage)
		return nil, errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	if len(defaultMethods) == 0 {
		errMessage := "could not load default methods with selector: %v"
		log.WithError(err).Debugf(errMessage, selector)
		return nil, errors.NotFoundError(errMessage, selector).ExtendComponent(component)
	}

	return defaultMethods, nil
}
