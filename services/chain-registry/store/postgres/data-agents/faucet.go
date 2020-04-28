package dataagents

import (
	"context"
	"fmt"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/go-pg/pg/v9"
	genuuid "github.com/satori/go.uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

const faucetComponentName = "chain-registry.store.faucet.pg"

type PGFaucetAgent struct {
	db *pg.DB
}

func NewPGFaucetAgent(
	db *pg.DB,
) *PGFaucetAgent {
	return &PGFaucetAgent{
		db: db,
	}
}

func (ag *PGFaucetAgent) RegisterFaucet(ctx context.Context, faucet *models.Faucet) error {
	if faucet.UUID == "" {
		faucet.UUID = genuuid.NewV4().String()
	}
	_, err := ag.db.ModelContext(ctx, faucet).Insert()
	if err != nil {
		errMessage := "Failed to register faucet"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(faucetComponentName)
	}

	return nil
}

func (ag *PGFaucetAgent) GetFaucets(ctx context.Context, filters map[string]string) ([]*models.Faucet, error) {
	faucets := make([]*models.Faucet, 0)

	req := ag.db.ModelContext(ctx, &faucets)
	for k, v := range filters {
		req.Where(fmt.Sprintf("%s = ?", k), v)
	}

	err := req.Select()
	if err != nil {
		errMessage := "Failed to get faucets"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.PostgresConnectionError(errMessage).ExtendComponent(faucetComponentName)
	}

	return faucets, nil
}

func (ag *PGFaucetAgent) GetFaucetsByTenant(ctx context.Context, filters map[string]string, tenantID string) ([]*models.Faucet, error) {
	faucets := make([]*models.Faucet, 0)

	req := ag.db.ModelContext(ctx, &faucets).
		Where("tenant_id = ?", tenantID)
	for k, v := range filters {
		req.Where(fmt.Sprintf("%s = ?", k), v)
	}

	err := req.Select()
	if err != nil {
		errMessage := "Failed to get faucets for tenant ID %s"
		log.FromContext(ctx).WithError(err).Error(errMessage, tenantID)
		return nil, errors.PostgresConnectionError(errMessage, tenantID).ExtendComponent(faucetComponentName)
	}

	return faucets, nil
}

func (ag *PGFaucetAgent) GetFaucetByUUID(ctx context.Context, uuid string) (*models.Faucet, error) {
	faucet := &models.Faucet{}

	err := ag.db.ModelContext(ctx, faucet).Where("uuid = ?", uuid).Select()
	if err != nil && err == pg.ErrNoRows {
		errMessage := "could not load faucet with chainUUID: %s"
		log.FromContext(ctx).WithError(err).Debugf(errMessage, uuid)
		return nil, errors.NotFoundError(errMessage, uuid).ExtendComponent(faucetComponentName)
	} else if err != nil {
		errMessage := "Failed to get faucet"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.PostgresConnectionError(errMessage).ExtendComponent(faucetComponentName)
	}

	return faucet, nil
}

func (ag *PGFaucetAgent) GetFaucetByUUIDAndTenant(ctx context.Context, uuid, tenantID string) (*models.Faucet, error) {
	faucet := &models.Faucet{}

	err := ag.db.ModelContext(ctx, faucet).Where("uuid = ?", uuid).Where("tenant_id = ?", tenantID).Select()
	if err != nil && err == pg.ErrNoRows {
		errMessage := "could not load faucet with chainUUID: %s and tenant: %s"
		log.FromContext(ctx).WithError(err).Debugf(errMessage, uuid, tenantID)
		return nil, errors.NotFoundError(errMessage, uuid, tenantID).ExtendComponent(faucetComponentName)
	} else if err != nil {
		errMessage := "Failed to get faucet from DB"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.PostgresConnectionError(errMessage).ExtendComponent(faucetComponentName)
	}

	return faucet, nil
}

func (ag *PGFaucetAgent) UpdateFaucetByUUID(ctx context.Context, uuid string, faucet *models.Faucet) error {
	res, err := ag.db.ModelContext(ctx, faucet).
		Where("uuid = ?", uuid).
		UpdateNotZero()

	if err != nil {
		errMessage := "Failed to update faucet"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(faucetComponentName)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		errMessage := "no faucet found with uuid %s"
		log.FromContext(ctx).WithError(err).Error(errMessage, faucet.UUID)
		return errors.NotFoundError(errMessage, faucet.UUID).ExtendComponent(faucetComponentName)
	}

	return nil
}

func (ag *PGFaucetAgent) DeleteFaucetByUUID(ctx context.Context, uuid string) error {
	faucet := &models.Faucet{}

	res, err := ag.db.ModelContext(ctx, faucet).Where("uuid = ?", uuid).Delete()
	if err != nil {
		errMessage := "Failed to delete faucet by UUID"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(faucetComponentName)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		errMessage := "no faucet found with chainUUID: %s"
		log.FromContext(ctx).WithError(err).Error(errMessage, uuid)
		return errors.NotFoundError(errMessage, uuid).ExtendComponent(faucetComponentName)
	}

	return nil
}

func (ag *PGFaucetAgent) DeleteFaucetByUUIDAndTenant(ctx context.Context, uuid, tenantID string) error {
	faucet := &models.Faucet{}

	res, err := ag.db.ModelContext(ctx, faucet).
		Where("uuid = ?", uuid).
		Where("tenant_id = ?", tenantID).Delete()
	if err != nil {
		errMessage := "Failed to delete faucet by UUID and tenant"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(faucetComponentName)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		errMessage := "no faucet found with uuid %s and tenant_id %s"
		log.FromContext(ctx).WithError(err).Error(errMessage, uuid, tenantID)
		return errors.NotFoundError(errMessage, uuid, tenantID).ExtendComponent(faucetComponentName)
	}

	return nil
}