package services

import (
	"context"
	"math/big"
)

// TxSender is used to send a transaction to the chain
type TxSender interface {
	// Send should send raw transaction
	Send(ctx context.Context, chainID *big.Int, raw string) error
}
