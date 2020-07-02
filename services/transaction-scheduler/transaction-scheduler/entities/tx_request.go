package entities

import (
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
)

type TxRequest struct {
	UUID           string
	IdempotencyKey string
	Schedule       *Schedule
	Params         *types.ETHTransactionParams
	Labels         map[string]string
	Annotations    *types.Annotations
	CreatedAt      time.Time
}
