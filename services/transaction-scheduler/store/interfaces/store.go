package interfaces

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
)

//go:generate mockgen -source=store.go -destination=mocks/store.go -package=mocks

type Store interface {
	Connect(ctx context.Context, conf interface{}) (DB, error)
}

type DB interface {
	Begin() (Tx, error)
	Agents
}

type Tx interface {
	database.Tx
	Agents
}

type Agents interface {
	TransactionRequest() TransactionRequestAgent
	Schedule() ScheduleAgent
	Job() JobAgent
	Log() LogAgent
}

// Interfaces data agents
type TransactionRequestAgent interface {
	SelectOrInsert(ctx context.Context, txRequest *models.TransactionRequest) error
	FindOneByIdempotencyKey(ctx context.Context, idempotencyKey string) (*models.TransactionRequest, error)
}

type ScheduleAgent interface {
	Insert(ctx context.Context, schedule *models.Schedule) error
	FindOneByUUID(ctx context.Context, scheduleUUID, tenantID string) (*models.Schedule, error)
}

type JobAgent interface {
	Insert(ctx context.Context, job *models.Job) error
	FindOneByUUID(ctx context.Context, jobUUID, tenantID string) (*models.Job, error)
}

type LogAgent interface {
	Insert(ctx context.Context, log *models.Log) error
}
