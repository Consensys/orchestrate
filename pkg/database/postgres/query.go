package postgres

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	healthz "github.com/heptiolabs/healthcheck"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
)

func Insert(ctx context.Context, db DB, models ...interface{}) error {
	_, err := db.ModelContext(ctx, models...).Insert()
	if err != nil {
		return handleError(err)
	}

	return nil
}

func InsertQuery(_ context.Context, q *orm.Query) error {
	_, err := q.Insert()
	if err != nil {
		return err
	}

	return nil
}

func Update(ctx context.Context, db DB, models ...interface{}) error {
	_, err := db.ModelContext(ctx, models...).WherePK().Update()
	if err != nil {
		return handleError(err)
	}
	return nil
}

func UpdateNotZero(ctx context.Context, q *orm.Query) error {
	_, err := q.Context(ctx).UpdateNotZero()
	if err != nil {
		return handleError(err)
	}
	return nil
}

func Delete(ctx context.Context, q *orm.Query) error {
	_, err := q.Context(ctx).Delete()
	if err != nil {
		return handleError(err)
	}

	return nil
}

func Select(ctx context.Context, q *orm.Query) error {
	err := q.Context(ctx).Select()
	if err != nil {
		return handleError(err)
	}

	return nil
}

func SelectColumn(ctx context.Context, q *orm.Query, result interface{}) error {
	err := q.Context(ctx).Select(result)
	if err != nil {
		return handleError(err)
	}

	return nil
}

func SelectOrInsert(ctx context.Context, q *orm.Query) error {
	_, err := q.Context(ctx).SelectOrInsert()
	if err != nil {
		return handleError(err)
	}

	return nil
}

func SelectOne(ctx context.Context, q *orm.Query) error {
	err := q.Context(ctx).First()
	if err != nil {
		return handleError(err)
	}
	return nil
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
