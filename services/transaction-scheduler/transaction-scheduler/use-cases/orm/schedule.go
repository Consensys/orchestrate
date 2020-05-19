package orm

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
)

func (orm *sORM) InsertSchedule(ctx context.Context, db store.DB, schedule *models.Schedule) error {
	return database.ExecuteInDBTx(db, func(tx database.Tx) error {
		dbtx := tx.(store.Tx)
		if schedule.ID != 0 {
			return errors.InvalidArgError("is not allowed to update a schedule")
		}

		for _, job := range schedule.Jobs {
			if job == nil {
				continue
			}
			if err := orm.InsertOrUpdateJob(ctx, dbtx, job); err != nil {
				return err
			}
		}

		if err := dbtx.Schedule().Insert(ctx, schedule); err != nil {
			return err
		}

		return nil
	})
}

func (orm *sORM) FetchScheduleByID(ctx context.Context, db store.DB, scheduleID int) (*models.Schedule, error) {
	schedule, err := db.Schedule().FindOneByID(ctx, scheduleID)
	if err != nil {
		return nil, err
	}

	for idx, job := range schedule.Jobs {
		schedule.Jobs[idx], err = db.Job().FindOneByUUID(ctx, job.UUID, schedule.TenantID)
		if err != nil {
			return schedule, err
		}
	}

	return schedule, nil
}

func (orm *sORM) FetchScheduleByUUID(ctx context.Context, db store.DB, scheduleUUID, tenantID string) (*models.Schedule, error) {
	schedule, err := db.Schedule().FindOneByUUID(ctx, scheduleUUID, tenantID)
	if err != nil {
		return nil, err
	}

	for idx, job := range schedule.Jobs {
		schedule.Jobs[idx], err = db.Job().FindOneByUUID(ctx, job.UUID, tenantID)
		if err != nil {
			return schedule, err
		}
	}

	return schedule, nil
}

func (orm *sORM) FetchAllSchedules(ctx context.Context, db store.DB, tenantID string) ([]*models.Schedule, error) {
	schedules, err := db.Schedule().FindAll(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	for idx, schedule := range schedules {
		for jdx, job := range schedule.Jobs {
			schedules[idx].Jobs[jdx], err = db.Job().FindOneByUUID(ctx, job.UUID, tenantID)
			if err != nil {
				return schedules, err
			}
		}
	}

	return schedules, nil
}
