package dataagents

import (
	"context"
	"time"

	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	pg "github.com/consensys/orchestrate/pkg/toolkit/database/postgres"
	"github.com/consensys/orchestrate/services/api/store"

	"github.com/consensys/orchestrate/services/api/store/models"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/types/entities"
	gopg "github.com/go-pg/pg/v9"
)

const accountDAComponent = "data-agents.account"

// PGAccount is an Account data agent for PostgreSQL
type PGAccount struct {
	db     pg.DB
	logger *log.Logger
}

// NewPGAccount creates a new PGAccount
func NewPGAccount(db pg.DB) store.AccountAgent {
	return &PGAccount{
		db:     db,
		logger: log.NewLogger().SetComponent(accountDAComponent),
	}
}

func (agent *PGAccount) Insert(ctx context.Context, account *models.Account) error {
	agent.db.ModelContext(ctx, account)
	err := pg.Insert(ctx, agent.db, account)
	if err != nil {
		errMsg := "failed to insert account"
		agent.logger.WithContext(ctx).WithError(err).Error(errMsg)
		return errors.FromError(err).SetMessage(errMsg).ExtendComponent(accountDAComponent)
	}

	return nil
}

// Insert Inserts a new job in DB
func (agent *PGAccount) Update(ctx context.Context, account *models.Account) error {

	if account.ID == 0 {
		errMsg := "cannot update account with missing ID"
		agent.logger.WithContext(ctx).Error(errMsg)
		return errors.InvalidArgError(errMsg).ExtendComponent(accountDAComponent)
	}

	account.UpdatedAt = time.Now().UTC()
	agent.db.ModelContext(ctx, account)
	err := pg.Update(ctx, agent.db, account)
	if err != nil {
		return errors.FromError(err).ExtendComponent(accountDAComponent)
	}

	return nil
}

func (agent *PGAccount) Search(ctx context.Context, filters *entities.AccountFilters, tenants []string, ownerID string) ([]*models.Account, error) {
	var accounts []*models.Account

	query := agent.db.ModelContext(ctx, &accounts)
	if len(filters.Aliases) > 0 {
		query = query.Where("alias in (?)", gopg.In(filters.Aliases))
	}
	if filters.TenantID != "" {
		query = query.Where("tenant_id = ?", filters.TenantID)
	}

	query = pg.WhereAllowedTenants(query, "tenant_id", tenants).Order("id ASC")
	query = pg.WhereAllowedOwner(query, "owner_id", ownerID)

	err := pg.Select(ctx, query)
	if err != nil {
		errMsg := "failed to search accounts"
		if !errors.IsNotFoundError(err) {
			agent.logger.WithContext(ctx).WithError(err).Error(errMsg)
		}
		return nil, errors.FromError(err).SetMessage(errMsg).ExtendComponent(accountDAComponent)
	}

	return accounts, nil
}

func (agent *PGAccount) FindOneByAddress(ctx context.Context, address string, tenants []string, ownerID string) (*models.Account, error) {
	account := &models.Account{}

	query := agent.db.ModelContext(ctx, account).Where("address = ?", address)

	query = pg.WhereAllowedTenants(query, "tenant_id", tenants)
	query = pg.WhereAllowedOwner(query, "owner_id", ownerID)

	err := pg.SelectOne(ctx, query)
	if err != nil {
		errMsg := "failed to find one account by address"
		if !errors.IsNotFoundError(err) {
			agent.logger.WithContext(ctx).WithError(err).Error(errMsg)
		}
		return nil, errors.FromError(err).SetMessage(errMsg).ExtendComponent(accountDAComponent)
	}

	return account, nil
}
