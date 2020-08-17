package entities

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
)

type JobSession struct {
	Job    *entities.Job
	Cancel context.CancelFunc
}
