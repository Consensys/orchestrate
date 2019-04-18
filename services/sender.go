package services

import (
	"context"
	"math/big"
)

// TxSender is used to send a transaction to the chain
type TxSender interface {
	// SendRawTransaction should send raw transaction
	SendRawTransaction(ctx context.Context, chainID *big.Int, raw string) error
}
