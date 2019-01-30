package ethereum

import (
	"context"
	"math/big"
)

// TxSenderEthClient is a minimal client interface required by a TxSender
type TxSenderEthClient interface {
	SendRawTransaction(ctx context.Context, chainID *big.Int, raw string) error
}

// TxSender is a sender that can manage one chain
type TxSender struct {
	ec TxSenderEthClient
}

// NewTxSender creates a new SingleChainSender
func NewTxSender(ec TxSenderEthClient) *TxSender {
	return &TxSender{ec}
}

// Send sends transaction
func (s *TxSender) Send(ctx context.Context, chainID *big.Int, raw string) error {
	return s.ec.SendRawTransaction(ctx, chainID, raw)
}
