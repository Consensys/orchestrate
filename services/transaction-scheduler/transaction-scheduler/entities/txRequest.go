package entities

import (
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
)

type TxRequest struct {
	IdempotencyKey string
	Schedule       *Schedule
	Params         *types.ETHTransactionParams
	Labels         map[string]string
	CreatedAt      time.Time
}
