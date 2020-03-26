package dataagents

import (
	"context"

	"github.com/go-pg/pg/v9/orm"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"

	"github.com/go-pg/pg/v9"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/models"
)

// PGTag is a tag data agent
type PGTag struct {
	db *pg.DB
}

// NewPGTag creates a new PGTag
func NewPGTag(db *pg.DB) *PGTag {
	return &PGTag{db: db}
}

// Insert Inserts a new tag in DB
func (agent *PGTag) Insert(ctx context.Context, tag *models.TagModel) error {
	tx := postgres.TxFromContext(ctx)
	if tx != nil {
		return agent.insert(tx.ModelContext(ctx, tag))
	}

	return agent.insert(agent.db.ModelContext(ctx, tag))
}

func (agent *PGTag) insert(query *orm.Query) error {
	_, err := query.
		OnConflict("ON CONSTRAINT tags_name_repository_id_key DO UPDATE").
		Set("artifact_id = ?artifact_id").
		Insert()
	if err != nil {
		errMessage := "could not create tag"
		log.WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	return nil
}

// FindAllByName Find all tags by name
func (agent *PGTag) FindAllByName(ctx context.Context, name string) ([]string, error) {
	var tags []string
	err := agent.db.ModelContext(ctx, (*models.TagModel)(nil)).
		Column("tag_model.name").
		Join("JOIN repositories AS registry ON registry.id = tag_model.repository_id").
		Where("registry.name = ?", name).
		OrderExpr("lower(tag_model.name)").
		Select(&tags)

	if err != nil {
		errMessage := "Failed to get tags from DB"
		log.WithError(err).Error(errMessage)
		return nil, errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	if len(tags) == 0 {
		errMessage := "could not load tags for name %s"
		log.WithError(err).Debugf(errMessage, name)
		return nil, errors.NotFoundError(errMessage, name).ExtendComponent(component)
	}

	return tags, nil
}
