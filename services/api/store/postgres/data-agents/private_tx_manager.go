package dataagents

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models"

	"github.com/gofrs/uuid"
	pg "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
)

const privateTxManagerDAComponent = "data-agents.private-tx-manager"

// PGPrivateTxManager is a Faucet data agent for PostgreSQL
type PGPrivateTxManager struct {
	db pg.DB
}

// NewPGPrivateTxManager creates a new PGPrivateTxManager
func NewPGPrivateTxManager(db pg.DB) store.PrivateTxManagerAgent {
	return &PGPrivateTxManager{db: db}
}

// Insert Inserts a new private transaction manager in DB
func (agent *PGPrivateTxManager) Insert(ctx context.Context, privateTxManager *models.PrivateTxManager) error {
	if privateTxManager.UUID == "" {
		privateTxManager.UUID = uuid.Must(uuid.NewV4()).String()
	}

	err := pg.Insert(ctx, agent.db, privateTxManager)
	if err != nil {
		return errors.FromError(err).ExtendComponent(privateTxManagerDAComponent)
	}

	return nil
}

func (agent *PGPrivateTxManager) Search(ctx context.Context, chainUUID string) ([]*models.PrivateTxManager, error) {
	var privateTxManagers []*models.PrivateTxManager

	query := agent.db.ModelContext(ctx, &privateTxManagers)
	if chainUUID != "" {
		query = query.Where("chain_uuid = ?", chainUUID)
	}

	err := pg.Select(ctx, query)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(privateTxManagerDAComponent)
	}

	return privateTxManagers, nil
}

func (agent *PGPrivateTxManager) Update(ctx context.Context, privateTxManager *models.PrivateTxManager) error {
	query := agent.db.ModelContext(ctx, privateTxManager).Where("uuid = ?", privateTxManager.UUID)

	err := pg.UpdateNotZero(ctx, query)
	if err != nil {
		return errors.FromError(err).ExtendComponent(privateTxManagerDAComponent)
	}

	return nil
}

func (agent *PGPrivateTxManager) Delete(ctx context.Context, privateTxManager *models.PrivateTxManager) error {
	query := agent.db.ModelContext(ctx, privateTxManager).Where("uuid = ?", privateTxManager.UUID)

	err := pg.Delete(ctx, query)
	if err != nil {
		return errors.FromError(err).ExtendComponent(privateTxManagerDAComponent)
	}

	return nil
}
