package dataagents

import (
	"context"

	pg "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models"
)

const repositoryDAComponent = "data-agents.repository"

// PGAccount is an Account data agent for PostgreSQL
type PGRepository struct {
	db pg.DB
}

func NewPGRepository(db pg.DB) store.RepositoryAgent {
	return &PGRepository{db: db}
}

func (agent *PGRepository) FindOne(ctx context.Context, name string) (*models.RepositoryModel, error) {
	model := &models.RepositoryModel{}
	query := agent.db.ModelContext(ctx, model).Where("LOWER(name) = LOWER(?)", name)
	err := pg.SelectOne(ctx, query)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(repositoryDAComponent)
	}

	return model, nil
}

func (agent *PGRepository) FindOneAndLock(ctx context.Context, name string) (*models.RepositoryModel, error) {
	model := &models.RepositoryModel{}
	query := agent.db.ModelContext(ctx, model).Where("LOWER(name) = LOWER(?)", name).For("UPDATE")
	err := pg.SelectOne(ctx, query)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(repositoryDAComponent)
	}

	return model, nil
}

func (agent *PGRepository) Insert(ctx context.Context, repository *models.RepositoryModel) error {
	err := pg.Insert(ctx, agent.db, repository)
	if err != nil {
		return errors.FromError(err).ExtendComponent(repositoryDAComponent)
	}

	return nil
}

func (agent *PGRepository) FindAll(ctx context.Context) ([]string, error) {
	var names []string
	query := agent.db.ModelContext(ctx, (*models.RepositoryModel)(nil)).
		Column("name").
		OrderExpr("lower(name)")

	err := pg.SelectColumn(ctx, query, &names)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(repositoryDAComponent)
	}

	return names, nil
}
