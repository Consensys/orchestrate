package services

import (
	"math/big"
)

// TxSender is used to send a transaction to the chain
type TxSender interface {
	// Send should send raw transaction
	Send(chainID *big.Int, raw string) error
}
