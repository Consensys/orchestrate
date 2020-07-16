package testutils

import (
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"

	"github.com/gofrs/uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
)

func FakeScheduleEntity() *entities.Schedule {
	scheduleUUID := uuid.Must(uuid.NewV4()).String()
	job := testutils.FakeJob()
	job.ScheduleUUID = scheduleUUID

	return &entities.Schedule{
		UUID: scheduleUUID,
		Jobs: []*types.Job{job},
	}
}

func FakeTxRequestEntity() *entities.TxRequest {
	return &entities.TxRequest{
		UUID:           uuid.Must(uuid.NewV4()).String(),
		Schedule:       FakeScheduleEntity(),
		IdempotencyKey: "IdempotencyKey",
		ChainName:      "chain",
		Params:         testutils.FakeETHTransactionParams(),
		CreatedAt:      time.Now(),
		Annotations:    &types.Annotations{},
	}
}

func FakeRawTxRequestEntity() *entities.TxRequest {
	return &entities.TxRequest{
		UUID:           uuid.Must(uuid.NewV4()).String(),
		Schedule:       FakeScheduleEntity(),
		IdempotencyKey: "IdempotencyKey",
		ChainName:      "chain",
		Params:         testutils.FakeRawTransactionParams(),
		CreatedAt:      time.Now(),
	}
}

func FakeTransferTxRequestEntity() *entities.TxRequest {
	return &entities.TxRequest{
		UUID:           uuid.Must(uuid.NewV4()).String(),
		Schedule:       FakeScheduleEntity(),
		IdempotencyKey: "IdempotencyKey",
		ChainName:      "chain",
		Params:         testutils.FakeTransferTransactionParams(),
		CreatedAt:      time.Now(),
	}
}

func FakeTesseraTxRequestEntity() *entities.TxRequest {
	return &entities.TxRequest{
		UUID:           uuid.Must(uuid.NewV4()).String(),
		Schedule:       FakeScheduleEntity(),
		IdempotencyKey: "IdempotencyKey",
		ChainName:      "chain",
		Params:         testutils.FakeTesseraTransactionParams(),
		CreatedAt:      time.Now(),
	}
}

func FakeOrionTxRequestEntity() *entities.TxRequest {
	return &entities.TxRequest{
		UUID:           uuid.Must(uuid.NewV4()).String(),
		Schedule:       FakeScheduleEntity(),
		IdempotencyKey: "IdempotencyKey",
		ChainName:      "chain",
		Params:         testutils.FakeOrionTransactionParams(),
		CreatedAt:      time.Now(),
	}
}
