package session

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/dynamic"
)

type Session interface {
	Run(ctx context.Context) error
}

type Builder interface {
	NewSession(node *dynamic.Node) (Session, error)
}
