package mock

import (
	"context"
	"math/big"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/types"
)

// Credit is a Mock crediting function
func Credit(ctx context.Context, r *types.Request) (*big.Int, error) {
	return r.Amount, nil
}

// Faucet is a mock Faucet
type Faucet struct{}

// NewFaucet creates a new mock faucet
func NewFaucet() *Faucet {
	return &Faucet{}
}

// Credit is mock crediting function
func (faucet *Faucet) Credit(ctx context.Context, r *types.Request) (*big.Int, error) {
	return Credit(ctx, r)
}
