package dataagents

import (
	"context"
	"fmt"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/go-pg/pg/v9"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

const chainComponentName = "chain-registry.store.chain.pg"

type PGChainAgent struct {
	db *pg.DB
}

func NewPGChainAgent(
	db *pg.DB,
) *PGChainAgent {
	return &PGChainAgent{
		db: db,
	}
}

func (ag *PGChainAgent) RegisterChain(ctx context.Context, chain *models.Chain) error {
	logger := log.FromContext(ctx)

	if err := chain.Validate(true); err != nil {
		logger.WithError(err).Errorf("could not register chain")
		return errors.DataError(err.Error())
	}

	_, err := ag.db.ModelContext(ctx, chain).Insert()
	if err != nil {
		logger.WithError(err).Errorf("could not register chain")
		if errors.IsAlreadyExistsError(err) {
			return errors.AlreadyExistsError("chain already exists").ExtendComponent(chainComponentName)
		}
		return errors.PostgresConnectionError("error registering chain").ExtendComponent(chainComponentName)
	}

	logger.WithFields(logrus.Fields{
		"uuid":      chain.UUID,
		"tenant.id": chain.TenantID,
		"name":      chain.Name,
		"urls":      chain.URLs,
	}).Infof("register chain")

	return nil
}

func (ag *PGChainAgent) GetChains(ctx context.Context, filters map[string]string) ([]*models.Chain, error) {
	var chains []*models.Chain

	req := ag.db.ModelContext(ctx, &chains)
	for k, v := range filters {
		req.Where(fmt.Sprintf("%s = ?", k), v)
	}

	err := req.Select()
	if err != nil {
		log.FromContext(ctx).WithError(err).Errorf("could not load chains")
		return nil, errors.PostgresConnectionError("error loading chains").ExtendComponent(chainComponentName)
	}

	return chains, nil
}

func (ag *PGChainAgent) GetChainsByTenant(ctx context.Context, filters map[string]string, tenantID string) ([]*models.Chain, error) {
	chains := make([]*models.Chain, 0)

	req := ag.db.ModelContext(ctx, &chains).
		Where("tenant_id = ?", tenantID)
	for k, v := range filters {
		req.Where(fmt.Sprintf("%s = ?", k), v)
	}

	err := req.Select()
	if err != nil {
		log.FromContext(ctx).
			WithField("tenant", tenantID).
			WithError(err).Errorf("could not load chains")
		return nil, errors.PostgresConnectionError("error loading chains for tenant %v", tenantID).ExtendComponent(chainComponentName)
	}

	return chains, nil
}

func (ag *PGChainAgent) GetChainByUUID(ctx context.Context, uuid string) (*models.Chain, error) {
	chain := &models.Chain{}

	err := ag.db.ModelContext(ctx, chain).Where("uuid = ?", uuid).Select()
	if err != nil && err == pg.ErrNoRows {
		return nil, errors.NotFoundError("chain %v does not exist", uuid).ExtendComponent(chainComponentName)
	} else if err != nil {
		log.FromContext(ctx).
			WithField("chain.uuid", uuid).
			WithError(err).Errorf("could not load chain")
		return nil, errors.PostgresConnectionError("error loading chain %v", uuid).ExtendComponent(chainComponentName)
	}

	return chain, nil
}

func (ag *PGChainAgent) GetChainByUUIDAndTenant(ctx context.Context, uuid, tenantID string) (*models.Chain, error) {
	chain := &models.Chain{}

	err := ag.db.ModelContext(ctx, chain).Where("uuid = ?", uuid).Where("tenant_id = ?", tenantID).Select()
	if err != nil && err == pg.ErrNoRows {
		return nil, errors.NotFoundError("chain %v does not exist in tenant %v", uuid, tenantID).ExtendComponent(chainComponentName)
	} else if err != nil {
		log.FromContext(ctx).
			WithField("chain.uuid", uuid).
			WithField("tenant", tenantID).
			WithError(err).Errorf("could not load chain")
		return nil, errors.PostgresConnectionError("error loading chain %v in tenant %v", uuid, tenantID).ExtendComponent(chainComponentName)
	}

	return chain, nil
}

func (ag *PGChainAgent) UpdateChainByName(ctx context.Context, chainName string, chain *models.Chain) error {
	logger := log.FromContext(ctx)
	if err := chain.Validate(false); err != nil {
		logger.WithError(err).Errorf("Failed to update chain by name")
		return errors.DataError(err.Error())
	}

	res, err := ag.db.ModelContext(ctx, chain).Where("name = ?", chainName).UpdateNotZero()
	if err != nil {
		errMessage := "Failed to update chain by name"
		logger.WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(chainComponentName)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {

		errMessage := "no chain found with tenant_id=%s and name=%s"
		logger.WithError(err).Error(errMessage, chain.TenantID, chain.Name)
		return errors.NotFoundError(errMessage, chain.TenantID, chain.Name).ExtendComponent(chainComponentName)
	}

	return nil
}

func (ag *PGChainAgent) UpdateChainByUUID(ctx context.Context, uuid string, chain *models.Chain) error {
	logger := log.FromContext(ctx)

	if err := chain.Validate(false); err != nil {
		logger.WithError(err).Errorf("Failed to update chain by UUID")
		return errors.DataError(err.Error())
	}

	res, err := ag.db.ModelContext(ctx, chain).Where("uuid = ?", uuid).UpdateNotZero()
	if err != nil {
		errMessage := "Failed to update chain by UUID"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(chainComponentName)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		errMessage := "no chain found with uuid %s for update"
		log.FromContext(ctx).WithError(err).Error(errMessage, chain.UUID)
		return errors.NotFoundError(errMessage, chain.UUID).ExtendComponent(chainComponentName)
	}

	return nil
}

func (ag *PGChainAgent) DeleteChainByUUID(ctx context.Context, uuid string) error {
	chain := &models.Chain{}

	res, err := ag.db.ModelContext(ctx, chain).Where("uuid = ?", uuid).Delete()
	if err != nil {
		errMessage := "Failed to delete chain by UUID"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(chainComponentName)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		errMessage := "no chain found with uuid %s for delete"
		log.FromContext(ctx).WithError(err).Error(errMessage, uuid)
		return errors.NotFoundError(errMessage, uuid).ExtendComponent(chainComponentName)
	}

	return nil
}

func (ag *PGChainAgent) DeleteChainByUUIDAndTenant(ctx context.Context, uuid, tenantID string) error {
	chain := &models.Chain{}

	res, err := ag.db.ModelContext(ctx, chain).Where("uuid = ?", uuid).Where("tenant_id = ?", tenantID).Delete()
	if err != nil {
		errMessage := "Failed to delete chain by UUID and tenant"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(chainComponentName)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		errMessage := "no chain found with uuid %s and tenant_id %s"
		log.FromContext(ctx).WithError(err).Error(errMessage, uuid, tenantID)
		return errors.NotFoundError(errMessage, uuid, tenantID).ExtendComponent(chainComponentName)
	}

	return nil
}

// // Insert Inserts a new contract in DB
// func (ag *PGChainAgent) Insert(
// 	ctx context.Context,
// 	chain *models.Chain,
// ) error {
// 	tx, err := agent.db.Begin()
// 	if err != nil {
// 		return errors.PostgresConnectionError("Failed to create DB transaction").ExtendComponent(chainComponentName)
// 	}
// 	pgctx := postgres.WithTx(ctx, tx)
//
// 	_, err = tx.ModelContext(pgctx, chain).
// 		Insert()
// 	if err != nil {
// 		errMessage := "could not create chain"
// 		log.WithError(err).Error(errMessage)
// 		return errors.PostgresConnectionError(errMessage).ExtendComponent(chainComponentName)
// 	}
//
// 	if chain.PrivateTxManagers != nil && len(chain.PrivateTxManagers) > 0 {
// 		err = agent.privateTxManangerDataAgent.InsertMultiple(pgctx, &chain.PrivateTxManagers)
// 		if err != nil {
// 			return errors.FromError(err).ExtendComponent(chainComponentName)
// 		}
// 	}
//
// 	return tx.Commit()
// }
