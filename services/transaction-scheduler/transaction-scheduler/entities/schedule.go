package entities

import (
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
)

type Schedule struct {
	UUID      string
	Jobs      []*types.Job
	CreatedAt time.Time
}
