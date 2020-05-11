package client

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
)

//go:generate mockgen -source=client.go -destination=mock/mock.go -package=mock

const component = "transaction-scheduler-client"

type TransactionSchedulerClient interface {
	SendTransaction(ctx context.Context, txRequest *types.TransactionRequest) (*types.TransactionResponse, error)
	GetSchedule(ctx context.Context, scheduleUUID string) (*types.ScheduleResponse, error)
}
