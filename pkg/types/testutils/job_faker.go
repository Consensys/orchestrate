package testutils

import (
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"

	"github.com/gofrs/uuid"
)

func FakeJob() *types.Job {
	return &types.Job{
		UUID:         uuid.Must(uuid.NewV4()).String(),
		ScheduleUUID: uuid.Must(uuid.NewV4()).String(),
		ChainUUID:    uuid.Must(uuid.NewV4()).String(),
		Type:         types.EthereumTransaction,
		Annotations:  &types.Annotations{ChainID: "888"},
		Logs:         []*types.Log{FakeLog()},
		CreatedAt:    time.Now(),
		Transaction:  FakeETHTransaction(),
	}
}
