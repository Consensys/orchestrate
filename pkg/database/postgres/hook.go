package postgres

import (
	"context"

	"github.com/go-pg/pg/v9"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
)

func handleError(err error) error {
	pgErr, ok := err.(pg.Error)
	if ok {
		switch {
		case pgErr.IntegrityViolation():
			return errors.AlreadyExistsError("integrity violation: item already exist - %v", pgErr)
		// List of codes could be found in https://www.postgresql.org/docs/10/errcodes-appendix.html
		case pgErr.Field('C')[0:2] == "08":
			return errors.PostgresConnectionError("database connection error - %v", pgErr)
		}
	}

	return errors.FromError(err)
}

type hook struct{}

func (h hook) BeforeQuery(ctx context.Context, q *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

func (h hook) AfterQuery(_ context.Context, q *pg.QueryEvent) error {
	log.Trace(q.FormattedQuery())
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
