package creditor

import (
	"context"
	"math/big"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/faucet"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/types"
)

// Controller is a controller that holds a list of account that should not be credited
type Controller struct{}

// NewController creates a new BlackList controller
func NewController() *Controller {
	return &Controller{}
}

// Control apply BlackList controller on a credit function
func (ctrl *Controller) Control(credit faucet.CreditFunc) faucet.CreditFunc {
	return func(ctx context.Context, r *types.Request) (*big.Int, error) {
		if len(r.FaucetsCandidates) == 0 {
			return nil, errors.FaucetWarning("no faucet candidates").ExtendComponent(component)
		}

		for key, candidate := range r.FaucetsCandidates {
			if candidate.Creditor.Hex() == r.Beneficiary.Hex() {
				delete(r.FaucetsCandidates, key)
			}
		}
		if len(r.FaucetsCandidates) == 0 {
			return nil, errors.FaucetSelfCreditWarning("attempt to credit the creditor").ExtendComponent(component)
		}

		return credit(ctx, r)
	}
}
