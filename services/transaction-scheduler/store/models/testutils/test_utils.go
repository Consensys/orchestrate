package testutils

import (
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"

	uuid "github.com/satori/go.uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
)

func FakeSchedule(tenantID string) *models.Schedule {
	if tenantID == "" {
		tenantID = multitenancy.DefaultTenantIDName
	}
	return &models.Schedule{
		TenantID: tenantID,
		UUID:     uuid.NewV4().String(),
		Jobs: []*models.Job{{
			UUID:        uuid.NewV4().String(),
			ChainUUID:   uuid.NewV4().String(),
			Type:        types.EthereumTransaction,
			Transaction: FakeTransaction(),
			Logs:        []*models.Log{{Status: types.StatusCreated, Message: "created message"}},
		}},
	}
}

func FakeTxRequest(scheduleID int) *models.TransactionRequest {
	fakeSchedule := FakeSchedule("")
	fakeSchedule.ID = scheduleID

	return &models.TransactionRequest{
		IdempotencyKey: uuid.NewV4().String(),
		RequestHash:    "requestHash",
		Params:         "{\"field0\": \"field0Value\"}",
		Schedules:      []*models.Schedule{fakeSchedule},
	}
}

func FakeTransaction() *models.Transaction {
	return &models.Transaction{
		UUID: uuid.NewV4().String(),
	}
}

func FakeJobModel(scheduleID int) *models.Job {
	job := &models.Job{
		UUID:      uuid.NewV4().String(),
		ChainUUID: uuid.NewV4().String(),
		Type:      types.EthereumTransaction,
		Schedule: &models.Schedule{
			ID:       scheduleID,
			TenantID: "_",
			UUID:     uuid.NewV4().String(),
		},
		Transaction: FakeTransaction(),
		Logs: []*models.Log{
			{UUID: uuid.NewV4().String(), Status: types.StatusCreated, Message: "created message"},
		},
		CreatedAt: time.Now(),
	}

	if scheduleID != 0 {
		job.ScheduleID = &scheduleID
	}

	return job
}

func FakeLog() *models.Log {
	return &models.Log{
		UUID:      uuid.NewV4().String(),
		Status:    types.StatusCreated,
		Job:       FakeJobModel(0),
		CreatedAt: time.Now(),
	}
}
