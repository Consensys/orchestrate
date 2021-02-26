package dataagents

import (
	"context"

	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
	"github.com/ConsenSys/orchestrate/services/api/store"
	"github.com/ConsenSys/orchestrate/services/api/store/models"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	pg "github.com/ConsenSys/orchestrate/pkg/toolkit/database/postgres"
	"github.com/gofrs/uuid"
)

const privateTxManagerDAComponent = "data-agents.private-tx-manager"

// PGPrivateTxManager is a Faucet data agent for PostgreSQL
type PGPrivateTxManager struct {
	db     pg.DB
	logger *log.Logger
}

// NewPGPrivateTxManager creates a new PGPrivateTxManager
func NewPGPrivateTxManager(db pg.DB) store.PrivateTxManagerAgent {
	return &PGPrivateTxManager{db: db, logger: log.NewLogger().SetComponent(privateTxManagerDAComponent)}
}

// Insert Inserts a new private transaction manager in DB
func (agent *PGPrivateTxManager) Insert(ctx context.Context, privateTxManager *models.PrivateTxManager) error {
	if privateTxManager.UUID == "" {
		privateTxManager.UUID = uuid.Must(uuid.NewV4()).String()
	}

	err := pg.Insert(ctx, agent.db, privateTxManager)
	if err != nil {
		agent.logger.WithContext(ctx).WithError(err).Error("failed to insert private tx manager")
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
		if !errors.IsNotFoundError(err) {
			agent.logger.WithContext(ctx).WithError(err).Error("failed to search private tx managers")
		}
		return nil, errors.FromError(err).ExtendComponent(privateTxManagerDAComponent)
	}

	return privateTxManagers, nil
}

func (agent *PGPrivateTxManager) Update(ctx context.Context, privateTxManager *models.PrivateTxManager) error {
	query := agent.db.ModelContext(ctx, privateTxManager).Where("uuid = ?", privateTxManager.UUID)

	err := pg.UpdateNotZero(ctx, query)
	if err != nil {
		agent.logger.WithContext(ctx).WithError(err).Error("failed to update private tx manager")
		return errors.FromError(err).ExtendComponent(privateTxManagerDAComponent)
	}

	return nil
}

func (agent *PGPrivateTxManager) Delete(ctx context.Context, privateTxManager *models.PrivateTxManager) error {
	query := agent.db.ModelContext(ctx, privateTxManager).Where("uuid = ?", privateTxManager.UUID)

	err := pg.Delete(ctx, query)
	if err != nil {
		agent.logger.WithContext(ctx).WithError(err).Error("failed to delete private tx manager")
		return errors.FromError(err).ExtendComponent(privateTxManagerDAComponent)
	}

	return nil
}
