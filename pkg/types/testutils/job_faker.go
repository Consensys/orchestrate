package testutils

import (
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"

	"github.com/gofrs/uuid"
)

func FakeJob() *entities.Job {
	return &entities.Job{
		UUID:         uuid.Must(uuid.NewV4()).String(),
		ScheduleUUID: uuid.Must(uuid.NewV4()).String(),
		ChainUUID:    uuid.Must(uuid.NewV4()).String(),
		TenantID:     utils.RandomString(6),
		Type:         utils.EthereumTransaction,
		InternalData: FakeInternalData(),
		Labels:       make(map[string]string),
		Logs:         []*entities.Log{FakeLog()},
		CreatedAt:    time.Now(),
		Status:       utils.StatusCreated,
		Transaction:  FakeETHTransaction(),
	}
}

func FakeInternalData() *entities.InternalData {
	return &entities.InternalData{
		ChainID:       "888",
		Priority:      utils.PriorityMedium,
		RetryInterval: 5 * time.Second,
	}
}
