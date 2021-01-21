package postgres

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	healthz "github.com/heptiolabs/healthcheck"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
	ierror "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/error"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
)

var alreadyExistErr = "entity already exists in DB"
var integrityErr = "insert integrity violation"

func Insert(ctx context.Context, db DB, models ...interface{}) *ierror.Error {
	_, err := db.ModelContext(ctx, models...).Insert()
	if err != nil {
		pgErr, ok := err.(pg.Error)
		if ok && errors.IsAlreadyExistsError(err) {
			return errors.AlreadyExistsError(alreadyExistErr)
		} else if ok && pgErr.IntegrityViolation() {
			return errors.ConstraintViolatedError(integrityErr)
		}

		return errors.PostgresConnectionError("error executing insert by model")
	}

	return nil
}

func InsertQuery(_ context.Context, q *orm.Query) *ierror.Error {
	_, err := q.Insert()
	if err != nil {
		pgErr, ok := err.(pg.Error)
		if ok && errors.IsAlreadyExistsError(err) {
			return errors.AlreadyExistsError(alreadyExistErr)
		} else if ok && pgErr.IntegrityViolation() {
			return errors.ConstraintViolatedError(integrityErr)
		}

		errMsg := "error executing insert by query"
		return errors.PostgresConnectionError(errMsg)
	}

	return nil
}

func Update(ctx context.Context, db DB, models ...interface{}) *ierror.Error {
	_, err := db.ModelContext(ctx, models...).WherePK().Update()
	if err != nil {
		pgErr, ok := err.(pg.Error)
		if ok && errors.IsAlreadyExistsError(err) {
			return errors.AlreadyExistsError("entity cannot be updated in DB")
		} else if ok && pgErr.IntegrityViolation() {
			return errors.PostgresConnectionError("update integrity violation")
		}

		return errors.PostgresConnectionError("error executing update")
	}
	return nil
}

func UpdateNotZero(ctx context.Context, q *orm.Query) *ierror.Error {
	_, err := q.Context(ctx).UpdateNotZero()
	if err != nil {
		pgErr, ok := err.(pg.Error)
		if ok && errors.IsAlreadyExistsError(err) {
			return errors.AlreadyExistsError("entity cannot be non zero updated in DB")
		} else if ok && pgErr.IntegrityViolation() {
			return errors.PostgresConnectionError("non zero update integrity violation")
		}

		return errors.PostgresConnectionError("error executing non zero update")
	}
	return nil
}

func Delete(ctx context.Context, q *orm.Query) *ierror.Error {
	_, err := q.Context(ctx).Delete()
	if err != nil {
		pgErr, ok := err.(pg.Error)
		if ok && pgErr.IntegrityViolation() {
			return errors.PostgresConnectionError("delete integrity violation")
		}

		return errors.PostgresConnectionError("error executing delete")
	}

	return nil
}

func Select(ctx context.Context, q *orm.Query) *ierror.Error {
	err := q.Context(ctx).Select()
	if err != nil && err == pg.ErrNoRows {
		return errors.NotFoundError("entities cannot be found")
	} else if err != nil {
		return errors.PostgresConnectionError("could not load entities")
	}

	return nil
}

func SelectColumn(ctx context.Context, q *orm.Query, result interface{}) *ierror.Error {
	err := q.Context(ctx).Select(result)
	if err != nil && err == pg.ErrNoRows {
		return errors.NotFoundError("entities cannot be found")
	} else if err != nil {
		return errors.PostgresConnectionError("could not load columns")
	}

	return nil
}

func SelectOrInsert(ctx context.Context, q *orm.Query) *ierror.Error {
	_, err := q.Context(ctx).SelectOrInsert()
	if err != nil {
		pgErr, ok := err.(pg.Error)
		if ok && pgErr.IntegrityViolation() {
			return errors.ConstraintViolatedError(integrityErr)
		}

		return errors.PostgresConnectionError("error executing select or insert")
	}

	return nil
}

func SelectOne(ctx context.Context, q *orm.Query) *ierror.Error {
	err := q.Context(ctx).First()
	if err != nil && err == pg.ErrNoRows {
		return errors.NotFoundError("entity does not exist")
	} else if err != nil {
		return errors.PostgresConnectionError("could not load entity")
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
