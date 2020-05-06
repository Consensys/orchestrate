package store

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
)

//go:generate mockgen -source=data-agents.go -destination=mocks/data-agents.go -package=mocks

type DataAgents struct {
	TransactionRequest TransactionRequestAgent
	ScheduleAgent      ScheduleAgent
	JobAgent           JobAgent
	LogAgent           LogAgent
}

// Interfaces data agents

type TransactionRequestAgent interface {
	SelectOrInsert(ctx context.Context, txRequest *models.TransactionRequest) error
	FindOneByIdempotencyKey(ctx context.Context, idempotencyKey string) (*models.TransactionRequest, error)
}

type ScheduleAgent interface {
	Insert(ctx context.Context, schedule *models.Schedule) error
}

type JobAgent interface {
	Insert(ctx context.Context, job *models.Job) error
	FindOneByUUID(ctx context.Context, jobUUID string) (*models.Job, error)
}

type LogAgent interface {
	Insert(ctx context.Context, log *models.Log) error
}
