package faucet

import (
	"context"
	"math/big"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/types"
)

// CreditFunc are functions expected to trigger an credit ether
type CreditFunc func(ctx context.Context, r *types.Request) (*big.Int, bool, error)

// Faucet is an interface for crediting an account with ether
type Faucet interface {
	// Credit should credit an account based on its own set of security rules
	// If credit is successful it should return amount credited and true
	// Credit should respond synchronously (not wait for a credit transaction to be mined)
	Credit(ctx context.Context, r *types.Request) (*big.Int, bool, error)
}
