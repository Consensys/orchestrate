package session

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/dynamic"
)

type Session interface {
	Run(ctx context.Context) error
}

type Builder interface {
	NewSession(chain *dynamic.Chain) (Session, error)
}

type SManager interface {
	Run(ctx context.Context) error
}
