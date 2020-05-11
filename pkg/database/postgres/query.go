package postgres

import (
	"context"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

func Insert(ctx context.Context, db orm.DB, models ...interface{}) error {
	_, err := db.ModelContext(ctx, models...).Insert()
	if err != nil {
		pgErr, ok := err.(pg.Error)
		if ok && errors.IsAlreadyExistsError(err) {
			return errors.AlreadyExistsError("entity already exists in DB")
		} else if ok && pgErr.IntegrityViolation() {
			return errors.PostgresConnectionError("integrity violation")
		}

		return errors.PostgresConnectionError("error executing insert")
	}
	return nil
}
