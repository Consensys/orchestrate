package testutils

import (
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"

	uuid "github.com/satori/go.uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
)

func FakeScheduleEntity() *entities.Schedule {
	scheduleUUID := uuid.NewV4().String()
	job := testutils.FakeJob()
	job.ScheduleUUID = scheduleUUID

	return &entities.Schedule{
		UUID: scheduleUUID,
		Jobs: []*types.Job{job},
	}
}

func FakeTxRequestEntity() *entities.TxRequest {
	return &entities.TxRequest{
		Schedule:       FakeScheduleEntity(),
		IdempotencyKey: "IdempotencyKey",
		Params:         testutils.FakeETHTransactionParams(),
		CreatedAt:      time.Now(),
	}
}

func FakeTesseraTxRequestEntity() *entities.TxRequest {
	return &entities.TxRequest{
		Schedule:       FakeScheduleEntity(),
		IdempotencyKey: "IdempotencyKey",
		Params:         testutils.FakeTesseraTransactionParams(),
		CreatedAt:      time.Now(),
	}
}

func FakeOrionTxRequestEntity() *entities.TxRequest {
	return &entities.TxRequest{
		Schedule:       FakeScheduleEntity(),
		IdempotencyKey: "IdempotencyKey",
		Params:         testutils.FakeOrionTransactionParams(),
		CreatedAt:      time.Now(),
	}
}
