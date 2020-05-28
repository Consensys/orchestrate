package testutils

import (
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"

	uuid "github.com/satori/go.uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
)

func FakeJobEntity() *entities.Job {
	return &entities.Job{
		UUID:         uuid.NewV4().String(),
		ScheduleUUID: uuid.NewV4().String(),
		Type:         tx.JobEthereumTransaction,
		Status:       entities.JobStatusCreated,
		CreatedAt:    time.Now(),
		Transaction:  FakeTransactionEntity(),
	}
}

func FakeScheduleEntity(chainUUID string) *entities.Schedule {
	scheduleUUID := uuid.NewV4().String()
	job := FakeJobEntity()
	job.ScheduleUUID = scheduleUUID

	return &entities.Schedule{
		UUID:      scheduleUUID,
		ChainUUID: chainUUID,
		Jobs:      []*entities.Job{job},
	}
}

func FakeTransactionEntity() *entities.Transaction {
	return &entities.Transaction{
		Hash:      "Hash",
		From:      "From",
		To:        "To",
		Nonce:     "Nonce",
		Value:     "Value",
		GasPrice:  "GasPrice",
		GasLimit:  "GasLimit",
		CreatedAt: time.Now(),
	}
}

func FakeTxRequestEntity() *entities.TxRequest {
	return &entities.TxRequest{
		Schedule:       FakeScheduleEntity("ChainUUID"),
		IdempotencyKey: "IdempotencyKey",
		Params:         FakeTxRequestParams(),
		CreatedAt:      time.Now(),
	}
}

func FakeTxRequestParams() *entities.TxRequestParams {
	return &entities.TxRequestParams{
		From:            &(&struct{ x string }{"From"}).x,
		To:              &(&struct{ x string }{"To"}).x,
		Value:           &(&struct{ x string }{"Value"}).x,
		GasPrice:        &(&struct{ x string }{"GasPrice"}).x,
		GasLimit:        &(&struct{ x string }{"GasLimit"}).x,
		MethodSignature: &(&struct{ x string }{"constructor(string,string)"}).x,
		Args:            []string{"val1", "val2"},
		Raw:             &(&struct{ x string }{"Raw"}).x,
		ContractName:    &(&struct{ x string }{"ContractName"}).x,
		ContractTag:     &(&struct{ x string }{"ContractTag"}).x,
	}
}
