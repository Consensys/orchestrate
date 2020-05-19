package orm

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
)

func (orm *sORM) InsertOrUpdateTransaction(ctx context.Context, db store.DB, tx *models.Transaction) error {
	if tx.ID == 0 {
		if err := db.Transaction().Insert(ctx, tx); err != nil {
			return err
		}
	} else {
		if err := db.Transaction().Update(ctx, tx); err != nil {
			return err
		}
	}

	return nil
}
