package testutils

import (
	uuid "github.com/satori/go.uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
)

func FakeSchedule() *models.Schedule {
	return &models.Schedule{
		TenantID: "tenantID",
		ChainID:  uuid.NewV4().String(),
		Jobs: []*models.Job{{
			UUID:        uuid.NewV4().String(),
			Type:        types.JobConstantinopleTransaction,
			Transaction: FakeTransaction(),
			Logs:        []*models.Log{{Status: types.LogStatusCreated, Message: "created message"}},
		}},
	}
}

func FakeTxRequest() *models.TransactionRequest {
	return &models.TransactionRequest{
		IdempotencyKey: uuid.NewV4().String(),
		RequestHash:    "requestHash",
		Params:         "{\"field0\": \"field0Value\"}",
		Schedule:       FakeSchedule(),
	}
}

func FakeTransaction() *models.Transaction {
	return &models.Transaction{
		UUID: uuid.NewV4().String(),
	}
}

func FakeJob(scheduleID int) *models.Job {
	return &models.Job{
		UUID:        uuid.NewV4().String(),
		Type:        types.JobConstantinopleTransaction,
		ScheduleID:  scheduleID,
		Transaction: FakeTransaction(),
		Logs:        []*models.Log{{Status: types.LogStatusCreated, Message: "created message"}},
	}
}

func FakeLog(jobID int) *models.Log {
	return &models.Log{
		UUID:   uuid.NewV4().String(),
		Status: types.LogStatusCreated,
		JobID:  jobID,
	}
}
