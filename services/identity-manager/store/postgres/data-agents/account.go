package dataagents

import (
	"context"

	gopg "github.com/go-pg/pg/v9"
	log "github.com/sirupsen/logrus"
	pg "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/identity-manager/store/models"
)

const accountDAComponent = "data-agents.account"

// NewPGAccount creates a new PGAccount
func NewPGAccount(db pg.DB) *PGAccount {
	return &PGAccount{db: db}
}

// PGAccount is an Account data agent for PostgreSQL
type PGAccount struct {
	db pg.DB
}

func (agent *PGAccount) Insert(ctx context.Context, account *models.Account) error {
	agent.db.ModelContext(ctx, account)
	err := pg.Insert(ctx, agent.db, account)
	if err != nil {
		return errors.FromError(err).ExtendComponent(accountDAComponent)
	}

	return nil
}

// Insert Inserts a new job in DB
func (agent *PGAccount) Update(ctx context.Context, account *models.Account) error {
	if account.ID == 0 {
		errMsg := "cannot update account with missing ID"
		log.WithContext(ctx).Error(errMsg)
		return errors.InvalidArgError(errMsg)
	}

	agent.db.ModelContext(ctx, account)
	err := pg.Update(ctx, agent.db, account)
	if err != nil {
		return errors.FromError(err).ExtendComponent(accountDAComponent)
	}

	return nil
}

func (agent *PGAccount) Search(ctx context.Context, filters *entities.AccountFilters, tenants []string) ([]*models.Account, error) {
	var idens []*models.Account

	query := agent.db.ModelContext(ctx, &idens)
	if len(filters.Aliases) > 0 {
		query = query.Where("alias in (?)", gopg.In(filters.Aliases))
	}

	query = pg.WhereAllowedTenants(query, "tenant_id", tenants).
		Order("id ASC")

	err := pg.Select(ctx, query)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(accountDAComponent)
	}

	return idens, nil
}

func (agent *PGAccount) FindOneByAddress(ctx context.Context, address string, tenants []string) (*models.Account, error) {
	var iden = &models.Account{}

	query := agent.db.ModelContext(ctx, iden).
		Where("address = ?", address)

	query = pg.WhereAllowedTenants(query, "tenant_id", tenants)

	err := pg.SelectOne(ctx, query)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(accountDAComponent)
	}

	return iden, nil
}
