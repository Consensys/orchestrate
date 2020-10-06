package postgres

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	healthz "github.com/heptiolabs/healthcheck"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	ierror "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/error"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

func Insert(ctx context.Context, db DB, models ...interface{}) *ierror.Error {
	logger := log.WithContext(ctx)
	_, err := db.ModelContext(ctx, models...).Insert()
	if err != nil {
		pgErr, ok := err.(pg.Error)
		if ok && errors.IsAlreadyExistsError(err) {
			errMsg := "entity already exists in DB"
			logger.WithError(err).Error(errMsg)
			return errors.AlreadyExistsError(errMsg)
		} else if ok && pgErr.IntegrityViolation() {
			errMsg := "insert integrity violation"
			logger.WithError(err).Error(errMsg)
			return errors.PostgresConnectionError(errMsg)
		}

		errMsg := "error executing insert"
		logger.WithError(err).Error(errMsg)
		return errors.PostgresConnectionError(errMsg)
	}

	return nil
}

func Update(ctx context.Context, db DB, models ...interface{}) *ierror.Error {
	logger := log.WithContext(ctx)
	_, err := db.ModelContext(ctx, models...).WherePK().Update()
	if err != nil {
		pgErr, ok := err.(pg.Error)
		if ok && errors.IsAlreadyExistsError(err) {
			errMsg := "entity cannot be updated in DB"
			logger.WithError(err).Error(errMsg)
			return errors.AlreadyExistsError(errMsg)
		} else if ok && pgErr.IntegrityViolation() {
			errMsg := "update integrity violation"
			logger.WithError(err).Error(errMsg)
			return errors.PostgresConnectionError(errMsg)
		}

		errMsg := "error executing update"
		logger.WithError(err).Error(errMsg)
		return errors.PostgresConnectionError(errMsg)
	}
	return nil
}

func Select(ctx context.Context, q *orm.Query) *ierror.Error {
	logger := log.WithContext(ctx)

	err := q.Context(ctx).Select()
	if err != nil && err == pg.ErrNoRows {
		return errors.NotFoundError("entities cannot be found")
	} else if err != nil {
		errMsg := "could not load entities"
		logger.WithError(err).Error(errMsg)
		return errors.PostgresConnectionError(errMsg)
	}

	return nil
}

func SelectOne(ctx context.Context, q *orm.Query) *ierror.Error {
	logger := log.WithContext(ctx)
	err := q.Context(ctx).First()
	if err != nil && err == pg.ErrNoRows {
		return errors.NotFoundError("entity does not exist")
	} else if err != nil {
		errMsg := "could not load entity"
		logger.WithError(err).Error(errMsg)
		return errors.PostgresConnectionError(errMsg)
	}
	return nil
}

func WhereFilters(query *orm.Query, filters map[string]string) *orm.Query {
	for k, v := range filters {
		query.Where(fmt.Sprintf("%s = ?", k), v)
	}
	return query
}

func WhereAllowedTenantsDefault(query *orm.Query, tenants []string) *orm.Query {
	return WhereAllowedTenants(query, "tenant_id", tenants)
}

func WhereAllowedTenants(query *orm.Query, field string, tenants []string) *orm.Query {
	if len(tenants) == 0 {
		return query
	}

	if utils.ContainsString(tenants, multitenancy.Wildcard) {
		return query
	}

	return query.Where(fmt.Sprintf("%s IN (?)", field), pg.In(tenants))
}

func Checker(db orm.DB) healthz.Check {
	return func() error {
		_, err := db.Exec("SELECT 1")
		return err
	}
}
