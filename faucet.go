package faucet

import (
	"context"
	"math/big"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
)

// CreditFunc are functions expected to trigger an credit ether
type CreditFunc func(ctx context.Context, r *services.FaucetRequest) (*big.Int, bool, error)

// ControlFunc are function expected to perform controls on a credit request
type ControlFunc func(f CreditFunc) CreditFunc

// CombineControls combines multiple controls into 1
func CombineControls(controls ...ControlFunc) ControlFunc {
	return func(f CreditFunc) CreditFunc {
		credit := f
		for i := len(controls); i > 0; i-- {
			credit = controls[i-1](credit)
		}
		return credit
	}
}

// ControlledFaucet a Faucet that credit only if some controls are valid
type ControlledFaucet struct {
	credit CreditFunc
}

// NewControlledFaucet create a new controlled faucet
func NewControlledFaucet(crediter CreditFunc, controls ...ControlFunc) *ControlledFaucet {
	return &ControlledFaucet{
		credit: CombineControls(controls...)(crediter),
	}
}

// Credit credit ethers
func (f *ControlledFaucet) Credit(ctx context.Context, r *services.FaucetRequest) (*big.Int, bool, error) {
	return f.credit(ctx, r)
}
