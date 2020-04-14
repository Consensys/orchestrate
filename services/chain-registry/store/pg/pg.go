package pg

import (
	"context"
	"fmt"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/go-pg/pg/v9"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

type Builder struct {
	postgres postgres.Manager
}

func NewBuilder(mngr postgres.Manager) *Builder {
	return &Builder{postgres: mngr}
}

func (b *Builder) Build(ctx context.Context, conf *Config) (*PG, error) {
	return New(b.postgres.Connect(ctx, conf.PG)), nil
}

// ChainRegistry is a traefik dynamic config registry based on PostgreSQL
type PG struct {
	db *pg.DB
}

// NewChainRegistry creates a new chain registry
func New(db *pg.DB) *PG {
	return &PG{db: db}
}

func (s *PG) RegisterChain(ctx context.Context, chain *types.Chain) error {
	logger := log.FromContext(ctx)

	if err := chain.Validate(true); err != nil {
		logger.WithError(err).Errorf("could not register chain")
		return errors.DataError(err.Error())
	}

	_, err := s.db.ModelContext(ctx, chain).Insert()
	if err != nil {
		logger.WithError(err).Errorf("could not register chain")
		if errors.IsAlreadyExistsError(err) {
			return errors.AlreadyExistsError("chain already exists").ExtendComponent(component)
		}
		return errors.PostgresConnectionError("error registering chain").ExtendComponent(component)
	}

	logger.WithFields(logrus.Fields{
		"uuid":      chain.UUID,
		"tenant.id": chain.TenantID,
		"name":      chain.Name,
		"urls":      chain.URLs,
	}).Infof("register chain")

	return nil
}

func (s *PG) GetChains(ctx context.Context, filters map[string]string) ([]*types.Chain, error) {
	var chains []*types.Chain

	req := s.db.ModelContext(ctx, &chains)
	for k, v := range filters {
		req.Where(fmt.Sprintf("%s = ?", k), v)
	}

	err := req.Select()
	if err != nil {
		log.FromContext(ctx).WithError(err).Errorf("could not load chains")
		return nil, errors.PostgresConnectionError("error loading chains").ExtendComponent(component)
	}

	return chains, nil
}

func (s *PG) GetChainsByTenant(ctx context.Context, filters map[string]string, tenantID string) ([]*types.Chain, error) {
	chains := make([]*types.Chain, 0)

	req := s.db.ModelContext(ctx, &chains).
		Where("tenant_id = ?", tenantID)
	for k, v := range filters {
		req.Where(fmt.Sprintf("%s = ?", k), v)
	}

	err := req.Select()
	if err != nil {
		log.FromContext(ctx).
			WithField("tenant", tenantID).
			WithError(err).Errorf("could not load chains")
		return nil, errors.PostgresConnectionError("error loading chains for tenant %v", tenantID).ExtendComponent(component)
	}

	return chains, nil
}

func (s *PG) GetChainByUUID(ctx context.Context, chainUUID string) (*types.Chain, error) {
	chain := &types.Chain{}

	err := s.db.ModelContext(ctx, chain).Where("uuid = ?", chainUUID).Select()
	if err != nil && err == pg.ErrNoRows {
		return nil, errors.NotFoundError("chain %v does not exist", chainUUID).ExtendComponent(component)
	} else if err != nil {
		log.FromContext(ctx).
			WithField("chain.uuid", chainUUID).
			WithError(err).Errorf("could not load chain")
		return nil, errors.PostgresConnectionError("error loading chain %v", chainUUID).ExtendComponent(component)
	}

	return chain, nil
}

func (s *PG) GetChainByUUIDAndTenant(ctx context.Context, chainUUID, tenantID string) (*types.Chain, error) {
	chain := &types.Chain{}

	err := s.db.ModelContext(ctx, chain).Where("uuid = ?", chainUUID).Where("tenant_id = ?", tenantID).Select()
	if err != nil && err == pg.ErrNoRows {
		return nil, errors.NotFoundError("chain %v does not exist in tenant %v", chainUUID, tenantID).ExtendComponent(component)
	} else if err != nil {
		log.FromContext(ctx).
			WithField("chain.uuid", chainUUID).
			WithField("tenant", tenantID).
			WithError(err).Errorf("could not load chain")
		return nil, errors.PostgresConnectionError("error loading chain %v in tenant %v", chainUUID, tenantID).ExtendComponent(component)
	}

	return chain, nil
}

func (s *PG) UpdateChainByName(ctx context.Context, chain *types.Chain) error {
	logger := log.FromContext(ctx)
	if err := chain.Validate(false); err != nil {
		logger.WithError(err).Errorf("Failed to update chain by name")
		return errors.DataError(err.Error())
	}

	res, err := s.db.ModelContext(ctx, chain).Where("name = ?", chain.Name).UpdateNotZero()
	if err != nil {
		errMessage := "Failed to update chain by name"
		logger.WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {

		errMessage := "no chain found with tenant_id=%s and name=%s"
		logger.WithError(err).Error(errMessage, chain.TenantID, chain.Name)
		return errors.NotFoundError(errMessage, chain.TenantID, chain.Name).ExtendComponent(component)
	}

	return nil
}

func (s *PG) UpdateChainByUUID(ctx context.Context, chain *types.Chain) error {
	logger := log.FromContext(ctx)

	if err := chain.Validate(false); err != nil {
		logger.WithError(err).Errorf("Failed to update chain by UUID")
		return errors.DataError(err.Error())
	}

	res, err := s.db.ModelContext(ctx, chain).WherePK().UpdateNotZero()
	if err != nil {
		errMessage := "Failed to update chain by UUID"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		errMessage := "no chain found with uuid %s for update"
		log.FromContext(ctx).WithError(err).Error(errMessage, chain.UUID)
		return errors.NotFoundError(errMessage, chain.UUID).ExtendComponent(component)
	}

	return nil
}

func (s *PG) DeleteChainByUUID(ctx context.Context, chainUUID string) error {
	chain := &types.Chain{}

	res, err := s.db.ModelContext(ctx, chain).Where("uuid = ?", chainUUID).Delete()
	if err != nil {
		errMessage := "Failed to delete chain by UUID"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		errMessage := "no chain found with uuid %s for delete"
		log.FromContext(ctx).WithError(err).Error(errMessage, chainUUID)
		return errors.NotFoundError(errMessage, chainUUID).ExtendComponent(component)
	}

	return nil
}

func (s *PG) DeleteChainByUUIDAndTenant(ctx context.Context, chainUUID, tenantID string) error {
	chain := &types.Chain{}

	res, err := s.db.ModelContext(ctx, chain).Where("uuid = ?", chainUUID).Where("tenant_id = ?", tenantID).Delete()
	if err != nil {
		errMessage := "Failed to delete chain by UUID and tenant"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		errMessage := "no chain found with uuid %s and tenant_id %s"
		log.FromContext(ctx).WithError(err).Error(errMessage, chainUUID, tenantID)
		return errors.NotFoundError(errMessage, chainUUID, tenantID).ExtendComponent(component)
	}

	return nil
}

func (s *PG) RegisterFaucet(ctx context.Context, faucet *types.Faucet) error {
	if faucet.UUID == "" {
		faucet.UUID = uuid.NewV4().String()
	}
	_, err := s.db.ModelContext(ctx, faucet).
		Insert()
	if err != nil {
		errMessage := "Failed to register faucet"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	return nil
}

func (s *PG) UpdateFaucetByUUID(ctx context.Context, faucet *types.Faucet) error {
	res, err := s.db.ModelContext(ctx, faucet).
		WherePK().
		UpdateNotZero()
	if err != nil {
		errMessage := "Failed to update faucet"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		errMessage := "no faucet found with uuid %s"
		log.FromContext(ctx).WithError(err).Error(errMessage, faucet.UUID)
		return errors.NotFoundError(errMessage, faucet.UUID).ExtendComponent(component)
	}

	return nil
}

func (s *PG) GetFaucets(ctx context.Context, filters map[string]string) ([]*types.Faucet, error) {
	faucets := make([]*types.Faucet, 0)

	req := s.db.ModelContext(ctx, &faucets)
	for k, v := range filters {
		req.Where(fmt.Sprintf("%s = ?", k), v)
	}

	err := req.Select()
	if err != nil {
		errMessage := "Failed to get faucets"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	return faucets, nil
}

func (s *PG) GetFaucetsByTenant(ctx context.Context, filters map[string]string, tenantID string) ([]*types.Faucet, error) {
	faucets := make([]*types.Faucet, 0)

	req := s.db.ModelContext(ctx, &faucets).
		Where("tenant_id = ?", tenantID)
	for k, v := range filters {
		req.Where(fmt.Sprintf("%s = ?", k), v)
	}

	err := req.Select()
	if err != nil {
		errMessage := "Failed to get faucets for tenant ID %s"
		log.FromContext(ctx).WithError(err).Error(errMessage, tenantID)
		return nil, errors.PostgresConnectionError(errMessage, tenantID).ExtendComponent(component)
	}

	return faucets, nil
}

func (s *PG) GetFaucetByUUID(ctx context.Context, chainUUID string) (*types.Faucet, error) {
	faucet := &types.Faucet{}

	err := s.db.ModelContext(ctx, faucet).Where("uuid = ?", chainUUID).Select()
	if err != nil && err == pg.ErrNoRows {
		errMessage := "could not load faucet with chainUUID: %s"
		log.FromContext(ctx).WithError(err).Debugf(errMessage, chainUUID)
		return nil, errors.NotFoundError(errMessage, chainUUID).ExtendComponent(component)
	} else if err != nil {
		errMessage := "Failed to get faucet"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	return faucet, nil
}

func (s *PG) GetFaucetByUUIDAndTenant(ctx context.Context, chainUUID, tenantID string) (*types.Faucet, error) {
	faucet := &types.Faucet{}

	err := s.db.ModelContext(ctx, faucet).Where("uuid = ?", chainUUID).Where("tenant_id = ?", tenantID).Select()
	if err != nil && err == pg.ErrNoRows {
		errMessage := "could not load faucet with chainUUID: %s and tenant: %s"
		log.FromContext(ctx).WithError(err).Debugf(errMessage, chainUUID, tenantID)
		return nil, errors.NotFoundError(errMessage, chainUUID, tenantID).ExtendComponent(component)
	} else if err != nil {
		errMessage := "Failed to get faucet from DB"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	return faucet, nil
}

func (s *PG) DeleteFaucetByUUID(ctx context.Context, chainUUID string) error {
	faucet := &types.Faucet{}

	res, err := s.db.ModelContext(ctx, faucet).Where("uuid = ?", chainUUID).Delete()
	if err != nil {
		errMessage := "Failed to delete faucet by UUID"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		errMessage := "no faucet found with chainUUID: %s"
		log.FromContext(ctx).WithError(err).Error(errMessage, chainUUID)
		return errors.NotFoundError(errMessage, chainUUID).ExtendComponent(component)
	}

	return nil
}

func (s *PG) DeleteFaucetByUUIDAndTenant(ctx context.Context, chainUUID, tenantID string) error {
	faucet := &types.Faucet{}

	res, err := s.db.ModelContext(ctx, faucet).Where("uuid = ?", chainUUID).Where("tenant_id = ?", tenantID).Delete()
	if err != nil {
		errMessage := "Failed to delete faucet by UUID and tenant"
		log.FromContext(ctx).WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		errMessage := "no faucet found with uuid %s and tenant_id %s"
		log.FromContext(ctx).WithError(err).Error(errMessage, chainUUID, tenantID)
		return errors.NotFoundError(errMessage, chainUUID, tenantID).ExtendComponent(component)
	}

	return nil
}
