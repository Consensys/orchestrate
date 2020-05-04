package dataagents

import (
	"context"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

func insert(ctx context.Context, query *orm.Query, component string) error {
	_, err := query.Insert()
	if err != nil {
		pgErr, ok := err.(pg.Error)
		if ok && pgErr.IntegrityViolation() || errors.IsAlreadyExistsError(err) {
			errMessage := "entity already exists in DB"
			log.WithContext(ctx).WithError(err).Error(errMessage)
			return errors.AlreadyExistsError(errMessage).ExtendComponent(component)
		}

		errMessage := "error executing insert"
		log.WithContext(ctx).WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(component)
	}

	return nil
}
