package entities

import (
	"time"
)

type TxRequest struct {
	UUID           string
	IdempotencyKey string
	ChainName      string
	Schedule       *Schedule
	Params         *ETHTransactionParams
	Labels         map[string]string
	InternalData   *InternalData
	CreatedAt      time.Time
}
