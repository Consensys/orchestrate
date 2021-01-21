package dataagents

import (
	"context"

	pg "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models"
)

const artifactDAComponent = "data-agents.artifact"

type PGArtifact struct {
	db     pg.DB
	logger *log.Logger
}

func NewPGArtifact(db pg.DB) store.ArtifactAgent {
	return &PGArtifact{
		db:     db,
		logger: log.NewLogger().SetComponent(artifactDAComponent),
	}
}

func (agent *PGArtifact) FindOneByABIAndCodeHash(ctx context.Context, abi, codeHash string) (*models.ArtifactModel, error) {
	model := &models.ArtifactModel{}
	query := agent.db.ModelContext(ctx, model).
		Where("abi = ?", abi).
		Where("codehash = ?", codeHash)

	err := pg.SelectOne(ctx, query)
	if err != nil {
		if !errors.IsNotFoundError(err) {
			agent.logger.WithContext(ctx).WithError(err).Error("failed to find artifact by ABI and CodeHash")
		}
		return nil, errors.FromError(err).ExtendComponent(repositoryDAComponent)
	}

	return model, nil
}

func (agent *PGArtifact) Insert(ctx context.Context, artifact *models.ArtifactModel) error {
	err := pg.Insert(ctx, agent.db, artifact)
	if err != nil {
		agent.logger.WithContext(ctx).WithError(err).Error("failed to insert artifact")
		return errors.FromError(err).ExtendComponent(artifactDAComponent)
	}

	return nil
}

func (agent *PGArtifact) SelectOrInsert(ctx context.Context, artifact *models.ArtifactModel) error {
	q := agent.db.ModelContext(ctx, artifact).Column("id").
		Where("abi = ?abi").Where("codehash = ?codehash").
		OnConflict("DO NOTHING").Returning("id")

	err := pg.SelectOrInsert(ctx, q)
	if err != nil {
		agent.logger.WithContext(ctx).WithError(err).Error("failed to select or insert artifact")
		return errors.FromError(err).ExtendComponent(artifactDAComponent)
	}

	return nil
}

func (agent *PGArtifact) FindOneByNameAndTag(ctx context.Context, name, tag string) (*models.ArtifactModel, error) {
	artifact := &models.ArtifactModel{}
	query := agent.db.ModelContext(ctx, artifact).
		Column("artifact_model.id", "abi", "bytecode", "deployed_bytecode").
		Join("JOIN tags AS t ON t.artifact_id = artifact_model.id").
		Join("JOIN repositories AS registry ON registry.id = t.repository_id").
		Where("LOWER(t.name) = LOWER(?)", tag).
		Where("LOWER(registry.name) = LOWER(?)", name)

	err := pg.SelectOne(ctx, query)
	if err != nil {
		agent.logger.WithContext(ctx).WithError(err).Error("failed to select artifact by name and tag")
		return nil, errors.FromError(err).ExtendComponent(artifactDAComponent)
	}

	return artifact, nil
}
