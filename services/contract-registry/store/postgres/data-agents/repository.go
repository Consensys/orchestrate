package dataagents

import (
	"context"

	"github.com/go-pg/pg/v9/orm"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres"

	"github.com/go-pg/pg/v9"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/contract-registry/store/models"
)

// PGRepository is a repository data agent
type PGRepository struct {
	db *pg.DB
}

// NewPGRepository creates a new PGRepository
func NewPGRepository(db *pg.DB) *PGRepository {
	return &PGRepository{db: db}
}

// SelectOrInsert Inserts a new repository in DB
func (agent *PGRepository) SelectOrInsert(ctx context.Context, repository *models.RepositoryModel) error {
	tx := postgres.TxFromContext(ctx)
	if tx != nil {
		return agent.selectOrInsert(tx.ModelContext(ctx, repository))
	}

	return agent.selectOrInsert(agent.db.ModelContext(ctx, repository))
}

func (agent *PGRepository) selectOrInsert(query *orm.Query) error {
	_, err := query.
		Column("id").
		Where("name = ?name").
		OnConflict("DO NOTHING").
		Returning("id").
		SelectOrInsert()
	if err != nil {
		errMessage := "could not create repository"
		log.WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	return nil
}

// FindAll Find all contracts
func (agent *PGRepository) FindAll(ctx context.Context) ([]string, error) {
	var names []string
	err := agent.db.ModelContext(ctx, (*models.RepositoryModel)(nil)).
		Column("name").
		OrderExpr("lower(name)").
		Select(&names)

	if err != nil {
		errMessage := "Failed to get catalog from DB"
		log.WithError(err).Error(errMessage)
		return nil, errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	return names, nil
}
