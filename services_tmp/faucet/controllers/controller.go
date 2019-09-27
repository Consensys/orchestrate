package controllers

import (
	"context"
	"math/big"

	"gitlab.com/ConsenSys/client/fr/core-stack/service/faucet.git/faucet"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/faucet.git/types"
)

// ControlFunc are function expected to perform controls on a credit request
type ControlFunc func(credit faucet.CreditFunc) faucet.CreditFunc

// CombineControls combines multiple controls into 1
func CombineControls(controls ...ControlFunc) ControlFunc {
	return func(credit faucet.CreditFunc) faucet.CreditFunc {
		for i := len(controls); i > 0; i-- {
			credit = controls[i-1](credit)
		}
		return credit
	}
}

// ControlledFaucet a Faucet that credits only if some controls are valid
type ControlledFaucet struct {
	credit faucet.CreditFunc
}

// NewControlledFaucet create a ControlledFaucet
func NewControlledFaucet(f faucet.Faucet, controls ...ControlFunc) *ControlledFaucet {
	return &ControlledFaucet{
		credit: CombineControls(controls...)(f.Credit),
	}
}

// Credit credit
func (f *ControlledFaucet) Credit(ctx context.Context, r *types.Request) (*big.Int, bool, error) {
	return f.credit(ctx, r)
}
