package orm

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
)

func (orm *sORM) InsertOrUpdateJob(ctx context.Context, db store.DB, job *models.Job) error {
	return database.ExecuteInDBTx(db, func(tx database.Tx) error {
		dbtx := tx.(store.Tx)
		if job.TransactionID != nil {
			if job.Transaction != nil && job.Transaction.ID != *job.TransactionID {
				err := errors.InvalidArgError("mismatched TransactionID")
				return err
			}
		} else if job.Transaction != nil {
			if err := orm.InsertOrUpdateTransaction(ctx, dbtx, job.Transaction); err != nil {
				return err
			}
			job.TransactionID = &job.Transaction.ID
		}

		if job.ScheduleID != nil {
			if job.Schedule != nil && job.Schedule.ID != *job.ScheduleID {
				err := errors.InvalidArgError("mismatched ScheduleID")
				return err
			}
		} else if job.Schedule != nil {
			// Schedule cannot be Update
			if job.Schedule.ID == 0 {
				if err := orm.InsertSchedule(ctx, dbtx, job.Schedule); err != nil {
					return err
				}
			}
			job.ScheduleID = &job.Schedule.ID
		}

		if job.ID == 0 {
			if err := dbtx.Job().Insert(ctx, job); err != nil {
				return err
			}
		} else {
			if err := dbtx.Job().Update(ctx, job); err != nil {
				return err
			}
		}

		for _, jobLog := range job.Logs {
			// Logs cannot be updated
			if jobLog == nil || jobLog.ID != 0 {
				continue
			}

			jobLog.JobID = &job.ID
			if err := orm.InsertLog(ctx, dbtx, jobLog); err != nil {
				return err
			}
		}

		return nil
	})
}
