package dataagents

import (
	"context"

	"github.com/go-pg/pg/v9"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/models"
)

// PGContract is a codehash data agent
type PGContract struct {
	db                  *pg.DB
	repositoryDataAgent store.RepositoryDataAgent
	artifactDataAgent   store.ArtifactDataAgent
	tagDataAgent        store.TagDataAgent
	methodDataAgent     store.MethodDataAgent
	eventDataAgent      store.EventDataAgent
}

// NewPGContract creates a new PGContract
func NewPGContract(
	db *pg.DB,
	repositoryDataAgent store.RepositoryDataAgent,
	artifactDataAgent store.ArtifactDataAgent,
	tagDataAgent store.TagDataAgent,
	methodDataAgent store.MethodDataAgent,
	eventDataAgent store.EventDataAgent,
) *PGContract {
	return &PGContract{
		db:                  db,
		repositoryDataAgent: repositoryDataAgent,
		artifactDataAgent:   artifactDataAgent,
		tagDataAgent:        tagDataAgent,
		methodDataAgent:     methodDataAgent,
		eventDataAgent:      eventDataAgent,
	}
}

// Insert Inserts a new contract in DB
func (agent *PGContract) Insert(
	ctx context.Context,
	name, tagName, abiRaw, bytecode, deployedBytecode, codeHash string,
	methods *[]*models.MethodModel,
	events *[]*models.EventModel,
) error {
	tx, err := agent.db.Begin()
	if err != nil {
		return errors.PostgresConnectionError("Failed to create DB transaction").ExtendComponent(component)
	}
	pgctx := postgres.WithTx(ctx, tx)

	repository := &models.RepositoryModel{
		Name: name,
	}
	err = agent.repositoryDataAgent.SelectOrInsert(pgctx, repository)
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	artifact := &models.ArtifactModel{
		Abi:              abiRaw,
		Bytecode:         bytecode,
		DeployedBytecode: deployedBytecode,
		Codehash:         codeHash,
	}
	err = agent.artifactDataAgent.SelectOrInsert(pgctx, artifact)
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	tag := &models.TagModel{
		Name:         tagName,
		RepositoryID: repository.ID,
		ArtifactID:   artifact.ID,
	}
	err = agent.tagDataAgent.Insert(pgctx, tag)
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	if methods != nil && len(*methods) > 0 {
		err = agent.methodDataAgent.InsertMultiple(pgctx, methods)
		if err != nil {
			return errors.FromError(err).ExtendComponent(component)
		}
	}

	if events != nil && len(*events) > 0 {
		err = agent.eventDataAgent.InsertMultiple(pgctx, events)
		if err != nil {
			return errors.FromError(err).ExtendComponent(component)
		}
	}

	return tx.Commit()
}
