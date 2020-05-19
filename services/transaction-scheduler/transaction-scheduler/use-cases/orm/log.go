package orm

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
)

func (orm *sORM) InsertLog(ctx context.Context, db store.DB, log *models.Log) error {
	if log.ID != 0 {
		return errors.InvalidArgError("is not allowed to update a job log")
	}

	if err := db.Log().Insert(ctx, log); err != nil {
		return err
	}

	return nil
}
