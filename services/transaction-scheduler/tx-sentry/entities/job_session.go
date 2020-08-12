package entities

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
)

type JobSession struct {
	Job    *types.Job
	Cancel context.CancelFunc
}
