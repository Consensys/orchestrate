package dataagents

import (
	"context"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/go-pg/pg/v9"
	"github.com/gofrs/uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
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
		faucet.UUID = uuid.Must(uuid.NewV4()).String()
	}
	_, err := ag.db.ModelContext(ctx, faucet).Insert()
	if err != nil {
		errMessage := "Failed to register faucet"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(faucetComponentName)
	}

	return nil
}

func (ag *PGFaucetAgent) GetFaucets(ctx context.Context, tenants []string, filters map[string]string) ([]*models.Faucet, error) {
	faucets := make([]*models.Faucet, 0)

	err := postgres.WhereFilters(
		postgres.WhereAllowedTenantsDefault(ag.db.ModelContext(ctx, &faucets), tenants),
		filters,
	).Select()
	if err != nil {
		errMessage := "Failed to get faucets"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.PostgresConnectionError(errMessage).ExtendComponent(faucetComponentName)
	}

	return faucets, nil
}

func (ag *PGFaucetAgent) GetFaucet(ctx context.Context, faucetUUID string, tenants []string) (*models.Faucet, error) {
	faucet := &models.Faucet{}

	err := postgres.WhereAllowedTenantsDefault(ag.db.ModelContext(ctx, faucet), tenants).
		Where("uuid = ?", faucetUUID).
		Select()
	if err != nil && err == pg.ErrNoRows {
		errMessage := "could not load faucet with chainUUID: %s"
		log.FromContext(ctx).WithError(err).Debugf(errMessage, faucetUUID)
		return nil, errors.NotFoundError(errMessage, faucetUUID).ExtendComponent(faucetComponentName)
	} else if err != nil {
		errMessage := "Failed to get faucet"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.PostgresConnectionError(errMessage).ExtendComponent(faucetComponentName)
	}

	return faucet, nil
}

func (ag *PGFaucetAgent) UpdateFaucet(ctx context.Context, faucetUUID string, tenants []string, faucet *models.Faucet) error {
	res, err := postgres.WhereAllowedTenantsDefault(ag.db.ModelContext(ctx, faucet), tenants).
		Where("uuid = ?", faucetUUID).
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

func (ag *PGFaucetAgent) DeleteFaucet(ctx context.Context, faucetUUID string, tenants []string) error {
	faucet := &models.Faucet{}

	res, err := postgres.WhereAllowedTenantsDefault(ag.db.ModelContext(ctx, faucet), tenants).
		Where("uuid = ?", faucetUUID).
		Delete()
	if err != nil {
		errMessage := "Failed to delete faucet by UUID"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(faucetComponentName)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		errMessage := "no faucet found with chainUUID: %s"
		log.FromContext(ctx).WithError(err).Error(errMessage, faucetUUID)
		return errors.NotFoundError(errMessage, faucetUUID).ExtendComponent(faucetComponentName)
	}

	return nil
}
