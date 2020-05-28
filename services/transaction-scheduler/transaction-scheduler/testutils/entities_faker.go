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

func FakeTransactionEntity() *entities.ETHTransaction {
	return &entities.ETHTransaction{
		From:           "From",
		To:             "To",
		Nonce:          "Nonce",
		Value:          "Value",
		GasPrice:       "GasPrice",
		GasLimit:       "GasLimit",
		Data:           "Data",
		Raw:            "Raw",
		PrivateFrom:    "PrivateFrom",
		PrivateFor:     []string{"val1", "val2"},
		PrivacyGroupID: "PrivacyGroupID",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
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
		From:            "From",
		To:              "To",
		Value:           "Value",
		GasPrice:        "GasPrice",
		GasLimit:        "GasLimit",
		MethodSignature: "constructor(string,string)",
		Args:            []string{"val1", "val2"},
		Raw:             "Raw",
		ContractName:    "ContractName",
		ContractTag:     "ContractTag",
		Nonce:           "1",
		PrivateFrom:     "PrivateFrom",
		PrivateFor:      []string{"val1", "val2"},
		PrivacyGroupID:  "PrivacyGroupID",
	}
}
