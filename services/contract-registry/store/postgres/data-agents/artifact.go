package dataagents

import (
	"context"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/models"
)

// PGArtifact is an artifact data agent
type PGArtifact struct {
	db *pg.DB
}

// NewPGArtifact creates a new PGArtifact
func NewPGArtifact(db *pg.DB) *PGArtifact {
	return &PGArtifact{db: db}
}

// SelectOrInsert Inserts a new artifact in DB
func (agent *PGArtifact) SelectOrInsert(ctx context.Context, artifact *models.ArtifactModel) error {
	tx := postgres.TxFromContext(ctx)
	if tx != nil {
		return agent.selectOrInsert(tx.ModelContext(ctx, artifact))
	}

	return agent.selectOrInsert(agent.db.ModelContext(ctx, artifact))
}

func (agent *PGArtifact) selectOrInsert(query *orm.Query) error {
	_, err := query.
		Column("id").
		Where("abi = ?abi").
		Where("codehash = ?codehash").
		OnConflict("DO NOTHING").
		Returning("id").
		SelectOrInsert()
	if err != nil {
		errMessage := "could not create artifact"
		log.WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	return nil
}

// SelectOne Selects the first found artifact in DB
func (agent *PGArtifact) FindOneByNameAndTag(ctx context.Context, name, tag string) (*models.ArtifactModel, error) {
	artifact := &models.ArtifactModel{}
	err := agent.db.ModelContext(ctx, artifact).
		Column("artifact_model.id", "abi", "bytecode", "deployed_bytecode").
		Join("JOIN tags AS t ON t.artifact_id = artifact_model.id").
		Join("JOIN repositories AS registry ON registry.id = t.repository_id").
		Where("t.name = ?", tag).
		Where("registry.name = ?", name).
		First()
	if err != nil && err == pg.ErrNoRows {
		errMessage := "could not load contract with name: %s and tag: %s"
		log.WithError(err).Debugf(errMessage, name, tag)
		return nil, errors.NotFoundError(errMessage, name, tag).ExtendComponent(component)
	} else if err != nil {
		errMessage := "Failed to get artifact from DB"
		log.WithError(err).Error(errMessage)
		return nil, errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	return artifact, nil
}
