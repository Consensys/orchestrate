package mock

import (
	"context"
	"math/big"

	"gitlab.com/ConsenSys/client/fr/core-stack/infra/faucet.git/faucet"
)

// Credit is a Mock crediting function
func Credit(ctx context.Context, r *faucet.Request) (*big.Int, bool, error) {
	return r.Amount, true, nil
}

// Faucet is a mock Faucet
type Faucet struct{}

// Credit is mock crediting function
func (faucet *Faucet) Credit(ctx context.Context, r *faucet.Request) (*big.Int, bool, error) {
	return Credit(ctx, r)
}
