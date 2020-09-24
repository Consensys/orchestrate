package testutils

import (
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"

	"github.com/gofrs/uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
)

func FakeSchedule(tenantID string) *models.Schedule {
	if tenantID == "" {
		tenantID = multitenancy.DefaultTenant
	}
	return &models.Schedule{
		TenantID: tenantID,
		UUID:     uuid.Must(uuid.NewV4()).String(),
		Jobs: []*models.Job{{
			UUID:        uuid.Must(uuid.NewV4()).String(),
			ChainUUID:   uuid.Must(uuid.NewV4()).String(),
			Type:        utils.EthereumTransaction,
			Transaction: FakeTransaction(),
			Logs:        []*models.Log{{Status: utils.StatusCreated, Message: "created message"}},
		}},
	}
}

func FakeTxRequest(scheduleID int) *models.TransactionRequest {
	fakeSchedule := FakeSchedule("")
	fakeSchedule.ID = scheduleID

	return &models.TransactionRequest{
		IdempotencyKey: utils.RandomString(16),
		ChainName:      "chain",
		RequestHash:    "requestHash",
		Params:         testutils.FakeETHTransactionParams(),
		Schedule:       fakeSchedule,
	}
}

func FakeTransaction() *models.Transaction {
	return &models.Transaction{
		UUID: uuid.Must(uuid.NewV4()).String(),
	}
}

func FakePrivateTx() *models.Transaction {
	tx := FakeTransaction()
	tx.PrivateFrom = "ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc="
	tx.PrivateFor = []string{"ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bd=", "ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Be="}
	return tx
}

func FakeJobModel(scheduleID int) *models.Job {
	job := &models.Job{
		UUID:      uuid.Must(uuid.NewV4()).String(),
		ChainUUID: uuid.Must(uuid.NewV4()).String(),
		Type:      utils.EthereumTransaction,
		Schedule: &models.Schedule{
			ID:       scheduleID,
			TenantID: "_",
			UUID:     uuid.Must(uuid.NewV4()).String(),
		},
		Transaction: FakeTransaction(),
		Logs: []*models.Log{
			{UUID: uuid.Must(uuid.NewV4()).String(), Status: utils.StatusCreated, Message: "created message", CreatedAt: time.Now()},
		},
		InternalData: &entities.InternalData{
			ChainID: "888",
		},
		CreatedAt: time.Now(),
		Labels:    make(map[string]string),
	}

	if scheduleID != 0 {
		job.ScheduleID = &scheduleID
	}

	return job
}

func FakeLog() *models.Log {
	return &models.Log{
		UUID:      uuid.Must(uuid.NewV4()).String(),
		Status:    utils.StatusCreated,
		Job:       FakeJobModel(0),
		CreatedAt: time.Now(),
	}
}
