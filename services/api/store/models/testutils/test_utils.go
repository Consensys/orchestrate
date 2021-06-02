package testutils

import (
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ConsenSys/orchestrate/pkg/types/entities"

	"github.com/ConsenSys/orchestrate/pkg/types/testutils"

	"github.com/ConsenSys/orchestrate/pkg/utils"

	"github.com/ConsenSys/orchestrate/pkg/multitenancy"
	"github.com/ConsenSys/orchestrate/services/api/store/models"
	"github.com/gofrs/uuid"
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
			Type:        entities.EthereumTransaction,
			Transaction: FakeTransaction(),
			Logs:        []*models.Log{{Status: entities.StatusCreated, Message: "created message"}},
		}},
	}
}

func FakeTxRequest(scheduleID int) *models.TransactionRequest {
	fakeSchedule := FakeSchedule("")
	fakeSchedule.ID = scheduleID

	return &models.TransactionRequest{
		IdempotencyKey: utils.RandString(16),
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
		Type:      entities.EthereumTransaction,
		Status:    entities.StatusCreated,
		Schedule: &models.Schedule{
			ID:       scheduleID,
			TenantID: "_",
			UUID:     uuid.Must(uuid.NewV4()).String(),
		},
		Transaction: FakeTransaction(),
		Logs: []*models.Log{
			{UUID: uuid.Must(uuid.NewV4()).String(), Status: entities.StatusCreated, Message: "created message", CreatedAt: time.Now()},
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
		Status:    entities.StatusCreated,
		Job:       FakeJobModel(0),
		CreatedAt: time.Now(),
	}
}

func FakeAccountModel() *models.Account {
	return &models.Account{
		Alias:               utils.RandString(10),
		TenantID:            "tenantID",
		Address:             ethcommon.HexToAddress(utils.RandHexString(12)).String(),
		PublicKey:           ethcommon.HexToHash(utils.RandHexString(12)).String(),
		CompressedPublicKey: ethcommon.HexToHash(utils.RandHexString(12)).String(),
		Attributes: map[string]string{
			"attr1": "val1",
			"attr2": "val2",
		},
	}
}

func FakeFaucetModel() *models.Faucet {
	return &models.Faucet{
		UUID:            uuid.Must(uuid.NewV4()).String(),
		TenantID:        "tenantID",
		Name:            "faucet-mainnet",
		ChainRule:       uuid.Must(uuid.NewV4()).String(),
		CreditorAccount: "0x5Cc634233E4a454d47aACd9fC68801482Fb02610",
		MaxBalance:      "60000000000000000",
		Amount:          "100000000000000000",
		Cooldown:        "10s",
	}
}

func FakeChainModel() *models.Chain {
	return &models.Chain{
		UUID:                      uuid.Must(uuid.NewV4()).String(),
		Name:                      "chain" + utils.RandString(5),
		TenantID:                  "tenantID",
		URLs:                      []string{"http://ganache:8545"},
		ChainID:                   "888",
		ListenerDepth:             0,
		ListenerCurrentBlock:      1,
		ListenerStartingBlock:     0,
		ListenerBackOffDuration:   "5s",
		ListenerExternalTxEnabled: false,
		PrivateTxManagers:         []*models.PrivateTxManager{},
		Labels: map[string]string{
			"label1": "value1",
		},
	}
}

func FakePrivateTxManagerModel() *models.PrivateTxManager {
	return &models.PrivateTxManager{
		UUID:      uuid.Must(uuid.NewV4()).String(),
		ChainUUID: "uuid.Must(uuid.NewV4()).String()",
		URL:       "http://tessera:8000",
		Type:      "Tessera",
		CreatedAt: time.Now(),
	}
}
