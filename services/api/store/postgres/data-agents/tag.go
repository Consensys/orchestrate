package dataagents

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"

	pg "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models"
)

const tagDAComponent = "data-agents.tag"

// PGAccount is an Account data agent for PostgreSQL
type PGTag struct {
	db     pg.DB
	logger *log.Logger
}

// NewPGAccount creates a new PGAccount
func NewPGTag(db pg.DB) store.TagAgent {
	return &PGTag{db: db, logger: log.NewLogger().SetComponent(tagDAComponent)}
}

func (agent *PGTag) Insert(ctx context.Context, tag *models.TagModel) error {
	query := agent.db.ModelContext(ctx, tag).
		OnConflict("ON CONSTRAINT tags_name_repository_id_key DO UPDATE").
		Set("artifact_id = ?artifact_id")

	err := pg.InsertQuery(ctx, query)
	if err != nil {
		agent.logger.WithContext(ctx).WithError(err).Error("failed to insert tag")
		return errors.FromError(err).ExtendComponent(tagDAComponent)
	}

	return nil
}
func (agent *PGTag) FindAllByName(ctx context.Context, name string) ([]string, error) {
	var tags []string
	query := agent.db.ModelContext(ctx, (*models.TagModel)(nil)).
		Column("tag_model.name").
		Join("JOIN repositories AS registry ON registry.id = tag_model.repository_id").
		Where("lower(registry.name) = lower(?)", name).
		OrderExpr("lower(tag_model.name)")

	err := pg.SelectColumn(ctx, query, &tags)
	if err != nil {
		if !errors.IsNotFoundError(err) {
			agent.logger.WithContext(ctx).WithError(err).Error("failed to fetch tag names")
		}
		return nil, errors.FromError(err).ExtendComponent(tagDAComponent)
	}

	return tags, nil
}
