package postgres

import (
	"context"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	ierror "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/error"
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

		return errors.PostgresConnectionError("error executing insert")
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

		return errors.PostgresConnectionError("error executing insert")
	}
	return nil
}

func Select(ctx context.Context, q *orm.Query) *ierror.Error {
	logger := log.WithContext(ctx)

	err := q.Context(ctx).Select()
	if err != nil && err == pg.ErrNoRows {
		errMsg := "entities cannot be found"
		logger.WithError(err).Error(errMsg)
		return errors.NotFoundError(errMsg)
	} else if err != nil {
		errMsg := "could not load entities"
		logger.WithError(err).Errorf(errMsg)
		return errors.PostgresConnectionError(errMsg)
	}

	return nil
}

func SelectOne(ctx context.Context, q *orm.Query) *ierror.Error {
	logger := log.WithContext(ctx)

	err := q.Context(ctx).First()
	if err != nil && err == pg.ErrNoRows {
		errMsg := "entity does not exist"
		logger.WithError(err).Error(errMsg)
		return errors.NotFoundError(errMsg)
	} else if err != nil {
		errMsg := "could not load entity"
		logger.WithError(err).Errorf(errMsg)
		return errors.PostgresConnectionError(errMsg)
	}

	return nil
}
