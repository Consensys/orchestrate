package orm

import (
	"context"

	storedb "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
)

//go:generate mockgen -source=orm.go -destination=./mocks/mock.go -package=mocks

type ORM interface {
	InsertOrUpdateJob(ctx context.Context, db storedb.DB, job *models.Job) error
	InsertLog(ctx context.Context, db storedb.DB, log *models.Log) error
	InsertSchedule(ctx context.Context, db storedb.DB, schedule *models.Schedule) error
	FetchScheduleByID(ctx context.Context, db storedb.DB, scheduleID int) (*models.Schedule, error)
	InsertOrUpdateTransaction(ctx context.Context, db storedb.DB, tx *models.Transaction) error
	FetchScheduleByUUID(ctx context.Context, db storedb.DB, scheduleUUID, tenantID string) (*models.Schedule, error)
	FetchAllSchedules(ctx context.Context, db storedb.DB, tenantID string) ([]*models.Schedule, error)
}

type sORM struct {
	ORM
}

func New() ORM {
	return &sORM{}
}
