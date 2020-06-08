package testutils

import (
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"

	uuid "github.com/satori/go.uuid"
)

func FakeJob() *types.Job {
	return &types.Job{
		UUID:         uuid.NewV4().String(),
		ScheduleUUID: uuid.NewV4().String(),
		ChainUUID:    uuid.NewV4().String(),
		Type:         types.EthereumTransaction,
		Logs:         []*types.Log{FakeLog()},
		CreatedAt:    time.Now(),
		Transaction:  FakeETHTransaction(),
	}
}
