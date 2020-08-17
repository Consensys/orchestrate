package testutils

import (
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"

	"github.com/gofrs/uuid"
)

func FakeTxRequest() *entities.TxRequest {
	return &entities.TxRequest{
		UUID:           uuid.Must(uuid.NewV4()).String(),
		Schedule:       FakeSchedule(),
		IdempotencyKey: "IdempotencyKey",
		ChainName:      "chain",
		Params:         FakeETHTransactionParams(),
		CreatedAt:      time.Now(),
		InternalData:   FakeInternalData(),
	}
}

func FakeRawTxRequest() *entities.TxRequest {
	return &entities.TxRequest{
		UUID:           uuid.Must(uuid.NewV4()).String(),
		Schedule:       FakeSchedule(),
		IdempotencyKey: "IdempotencyKey",
		ChainName:      "chain",
		Params:         FakeRawTransactionParams(),
		CreatedAt:      time.Now(),
		InternalData:   FakeInternalData(),
	}
}

func FakeTransferTxRequest() *entities.TxRequest {
	return &entities.TxRequest{
		UUID:           uuid.Must(uuid.NewV4()).String(),
		Schedule:       FakeSchedule(),
		IdempotencyKey: "IdempotencyKey",
		ChainName:      "chain",
		Params:         FakeTransferTransactionParams(),
		CreatedAt:      time.Now(),
		InternalData:   FakeInternalData(),
	}
}

func FakeTesseraTxRequest() *entities.TxRequest {
	return &entities.TxRequest{
		UUID:           uuid.Must(uuid.NewV4()).String(),
		Schedule:       FakeSchedule(),
		IdempotencyKey: "IdempotencyKey",
		ChainName:      "chain",
		Params:         FakeTesseraTransactionParams(),
		CreatedAt:      time.Now(),
		InternalData:   FakeInternalData(),
	}
}

func FakeOrionTxRequest() *entities.TxRequest {
	return &entities.TxRequest{
		UUID:           uuid.Must(uuid.NewV4()).String(),
		Schedule:       FakeSchedule(),
		IdempotencyKey: "IdempotencyKey",
		ChainName:      "chain",
		Params:         FakeOrionTransactionParams(),
		CreatedAt:      time.Now(),
		InternalData:   FakeInternalData(),
	}
}
