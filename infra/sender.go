package infra

import (
	"context"
	"math/big"
)

// TxSender is used to send a transaction to the chain
type TxSender interface {
	// Send should send raw transaction
	Send(chainID *big.Int, raw string) error
}

// SimpleSender is a sender that can manage one chain
type SimpleSender struct {
	ec *EthClient
}

// NewSimpleSender creates a new SimpleSender
func NewSimpleSender(ec *EthClient) *SimpleSender {
	return &SimpleSender{ec}
}

// Send sends transaction
func (s *SimpleSender) Send(chainID *big.Int, raw string) error {
	return s.ec.SendRawTransaction(context.Background(), raw)
}
