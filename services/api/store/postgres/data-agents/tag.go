package dataagents

import (
	"context"
	"strings"

	"github.com/consensys/orchestrate/pkg/toolkit/app/log"

	"github.com/consensys/orchestrate/pkg/errors"
	pg "github.com/consensys/orchestrate/pkg/toolkit/database/postgres"
	"github.com/consensys/orchestrate/services/api/store"
	"github.com/consensys/orchestrate/services/api/store/models"
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
		Join("JOIN repositories AS repo ON repo.id = tag_model.repository_id").
		Where("lower(repo.name) = ?", strings.ToLower(name)).
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
