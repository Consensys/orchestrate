package postgres

import (
	"context"
	"fmt"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/go-pg/pg/v9"
	log "github.com/sirupsen/logrus"
)

func handleError(err error) error {
	if pg.ErrNoRows == err {
		return errors.NotFoundError("data cannot be found")
	}
	if pg.ErrMultiRows == err {
		return errors.DataCorruptedError("multiple rows found, only expected one")
	}

	pgErr, ok := err.(pg.Error)
	if ok {
		switch {
		case pgErr.IntegrityViolation():
			return errors.ConstraintViolatedError("database integrity violation")
		case pgErr.Field('C')[0:2] == "22":
			return errors.InvalidParameterError("database input data").AppendReason(pgErr.Error())
		// List of codes could be found in https://www.postgresql.org/docs/10/errcodes-appendix.html
		case pgErr.Field('C')[0:2] == "08":
			return errors.PostgresConnectionError("database connection error").AppendReason(pgErr.Error())
		}
	}

	return fmt.Errorf("database internal error - %s", err.Error())
}

type hook struct{}

func (h hook) BeforeQuery(ctx context.Context, q *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

func (h hook) AfterQuery(ctx context.Context, q *pg.QueryEvent) error {
	log.WithContext(ctx).Trace(q.FormattedQuery())
	if q.Err != nil {
		q.Err = handleError(q.Err)
		return q.Err
	}
	return nil
}

func New(opts *pg.Options) *pg.DB {
	db := pg.Connect(opts)
	db.AddQueryHook(hook{})
	return db
}
