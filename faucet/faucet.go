package faucet

import (
	"context"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

// Request holds information for a Faucet Credit Request
type Request struct {
	ChainID     *big.Int
	Creditor    ethcommon.Address
	Beneficiary ethcommon.Address
	Amount      *big.Int
}

// CreditFunc are functions expected to trigger an credit ether
type CreditFunc func(ctx context.Context, r *Request) (*big.Int, bool, error)

// Faucet is an interface for crediting an account with ether
type Faucet interface {
	// Credit should credit an account based on its own set of security rules
	// If credit is successful it should return amount credited and true
	// Credit should respond synchronously (not wait for a credit transaction to be mined)
	Credit(ctx context.Context, r *Request) (*big.Int, bool, error)
}
