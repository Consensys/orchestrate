package handlers

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
)

// TxSender is used to send a transaction to the chain
type TxSender interface {
	// Send should send raw transaction
	Send(chainID *big.Int, raw string) error
}

// SimpleSender is a sender that can manage one chain
type SimpleSender struct {
	ec *infra.EthClient
}

// NewSimpleSender creates a new SimpleSender
func NewSimpleSender(ec *infra.EthClient) *SimpleSender {
	return &SimpleSender{ec}
}

// Send sends transaction
func (s *SimpleSender) Send(chainID *big.Int, raw string) error {
	return s.ec.SendRawTransaction(context.Background(), raw)
}

// Sender creates a Sender handler
func Sender(sender TxSender) infra.HandlerFunc {
	return func(ctx *infra.Context) {
		if len(ctx.T.Tx().Raw()) == 0 {
			// Tx is not ready
			// TODO: handle case
			ctx.Abort()
			return
		}

		err := sender.Send(ctx.T.Chain().ID, hexutil.Encode(ctx.T.Tx().Raw()))
		if err != nil {
			// TODO: handle error
			ctx.AbortWithError(err)
			return
		}
	}
}
