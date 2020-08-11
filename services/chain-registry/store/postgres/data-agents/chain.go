package dataagents

import (
	"context"
	"fmt"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/go-pg/pg/v9"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
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

// TODO: RegisterChain should be an atomic database transaction
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
		"chainUUID": chain.UUID,
		"tenantID":  chain.TenantID,
		"name":      chain.Name,
		"urls":      chain.URLs,
	}).Infof("registered chain")

	for _, prixTxManager := range chain.PrivateTxManagers {
		_, err := ag.db.ModelContext(ctx, prixTxManager).Insert()
		if err != nil {
			logger.WithError(err).Errorf("could not register priv tx manager")
			return errors.PostgresConnectionError("error registering priv tx manager").ExtendComponent(chainComponentName)
		}

		logger.WithFields(logrus.Fields{
			"uuid":      prixTxManager.UUID,
			"chainUUID": chain.UUID,
			"type":      prixTxManager.Type,
			"url":       prixTxManager.URL,
		}).Infof("register private tx manager")
	}

	return nil
}

func (ag *PGChainAgent) GetChains(ctx context.Context, tenants []string, filters map[string]string) ([]*models.Chain, error) {
	var chains []*models.Chain

	err := postgres.WhereFilters(
		postgres.WhereAllowedTenantsDefault(ag.db.ModelContext(ctx, &chains), tenants),
		filters,
	).Select()
	if err != nil {
		log.FromContext(ctx).WithError(err).Errorf("could not load chains")
		return nil, errors.PostgresConnectionError("error loading chains").ExtendComponent(chainComponentName)
	}

	for idx, c := range chains {
		err = ag.db.ModelContext(ctx, &chains[idx].PrivateTxManagers).
			Where(fmt.Sprintf(`chain_uuid = '%s'`, c.UUID)).
			Select()

		if err != nil {
			log.FromContext(ctx).WithError(err).Errorf("could not load private tx managers")
			return nil, errors.PostgresConnectionError("error loading chains").ExtendComponent(chainComponentName)
		}
	}

	return chains, nil
}

func (ag *PGChainAgent) GetChain(ctx context.Context, uuid string, tenants []string) (*models.Chain, error) {
	chain := &models.Chain{}

	err := postgres.WhereAllowedTenantsDefault(ag.db.ModelContext(ctx, chain), tenants).
		Where("uuid = ?", uuid).
		Select()
	if err != nil && err == pg.ErrNoRows {
		return nil, errors.NotFoundError("chain %v does not exist", uuid).ExtendComponent(chainComponentName)
	} else if err != nil {
		log.FromContext(ctx).
			WithField("chainUUID", uuid).
			WithError(err).Errorf("could not load chain")
		return nil, errors.PostgresConnectionError("error loading chain %v", uuid).ExtendComponent(chainComponentName)
	}

	err = ag.db.ModelContext(ctx, &chain.PrivateTxManagers).
		Where(fmt.Sprintf(`chain_uuid = '%s'`, chain.UUID)).
		Select()
	if err != nil {
		log.FromContext(ctx).WithError(err).Errorf("could not load private tx managers")
		return nil, errors.PostgresConnectionError("error loading chain %v", uuid).ExtendComponent(chainComponentName)
	}

	return chain, nil
}

func (ag *PGChainAgent) UpdateChainByName(ctx context.Context, chainName string, tenants []string, chain *models.Chain) error {
	logger := log.FromContext(ctx)
	if err := chain.Validate(false); err != nil {
		logger.WithError(err).Errorf("Failed to update chain by name")
		return errors.DataError(err.Error())
	}

	res, err := postgres.WhereAllowedTenantsDefault(
		ag.db.ModelContext(ctx, chain).Where("name = ?", chainName), tenants,
	).UpdateNotZero()

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

	if len(chain.PrivateTxManagers) > 0 {
		err = ag.updateChainPrivateTxManagers(ctx, chain.UUID, chain.PrivateTxManagers)
		if err != nil {
			logger.WithError(err).Error("cannot delete private tx managers")
			return errors.NotFoundError("cannot delete private tx managers").ExtendComponent(chainComponentName)
		}
	}

	return nil
}

func (ag *PGChainAgent) UpdateChain(ctx context.Context, uuid string, tenants []string, chain *models.Chain) error {
	logger := log.FromContext(ctx)

	if err := chain.Validate(false); err != nil {
		logger.WithError(err).Errorf("Failed to update chain by UUID")
		return errors.DataError(err.Error())
	}

	res, err := postgres.WhereAllowedTenantsDefault(ag.db.ModelContext(ctx, chain), tenants).
		Where("uuid = ?", uuid).
		UpdateNotZero()
	if err != nil {
		errMessage := "Failed to update chain by UUID"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(chainComponentName)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		errMessage := "no chain found with uuid %s for update"
		log.FromContext(ctx).WithError(err).Errorf(errMessage, uuid)
		return errors.NotFoundError(errMessage, uuid).ExtendComponent(chainComponentName)
	}

	if len(chain.PrivateTxManagers) > 0 {
		err = ag.updateChainPrivateTxManagers(ctx, uuid, chain.PrivateTxManagers)
		if err != nil {
			logger.WithError(err).Error("cannot delete private tx managers")
			return errors.DataError("cannot delete private tx managers").ExtendComponent(chainComponentName)
		}
	}

	return nil
}

func (ag *PGChainAgent) DeleteChain(ctx context.Context, uuid string, tenants []string) error {
	chain := &models.Chain{}

	res, err := postgres.WhereAllowedTenantsDefault(ag.db.ModelContext(ctx, chain), tenants).
		Where("uuid = ?", uuid).
		Delete()
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

func (ag *PGChainAgent) updateChainPrivateTxManagers(ctx context.Context, chainUUID string, privateTxManagers []*models.PrivateTxManagerModel) error {
	logger := log.FromContext(ctx)

	// Remove previous PrivateTxManagers to perform an full update
	_, err := ag.db.ModelContext(ctx, &models.PrivateTxManagerModel{}).
		Where("chain_uuid = ?", chainUUID).
		Delete()

	if err != nil {
		return err
	}

	logger.WithFields(logrus.Fields{
		"chainUUID": chainUUID,
	}).Infof("removed privateTxManagers")

	for _, prixTxManager := range privateTxManagers {
		prixTxManager.ChainUUID = chainUUID
		_, err := ag.db.ModelContext(ctx, prixTxManager).Insert()
		if err != nil {
			return err
		}

		logger.WithFields(logrus.Fields{
			"uuid":      prixTxManager.UUID,
			"chainUUID": chainUUID,
			"type":      prixTxManager.Type,
			"url":       prixTxManager.URL,
		}).Infof("registered private tx manager")
	}

	return nil
}
