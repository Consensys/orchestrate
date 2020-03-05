package maxbalance

import (
	"context"
	"math/big"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/faucet"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/types"
)

// Controller is a controller that ensures an address can not be credit above a given limit
type Controller struct {
	ChainStateReader ethclient.ChainStateReader
}

// NewController creates a new max balance controller
func NewController(chainStateReader ethclient.ChainStateReader) *Controller {
	return &Controller{
		chainStateReader,
	}
}

// Control apply MaxBalance controller on a credit function
func (ctrl *Controller) Control(credit faucet.CreditFunc) faucet.CreditFunc {
	return func(ctx context.Context, r *types.Request) (*big.Int, error) {
		if len(r.FaucetsCandidates) == 0 {
			return nil, errors.FaucetWarning("no faucet candidates").ExtendComponent(component)
		}

		// Retrieve account balance
		balance, err := ctrl.ChainStateReader.BalanceAt(ctx, r.ChainURL, r.Beneficiary, nil)
		if err != nil {
			return nil, errors.FromError(err).ExtendComponent(component)
		}

		// Ensure MaxBalance is respected
		for key, candidate := range r.FaucetsCandidates {
			if balance.Add(balance, candidate.Amount).Cmp(candidate.MaxBalance) > 0 {
				delete(r.FaucetsCandidates, key)
			}
		}
		if len(r.FaucetsCandidates) == 0 {
			// Do not credit if final balance would be superior to max authorized
			return nil, errors.FaucetWarning("account balance too high").ExtendComponent(component)
		}

		return credit(ctx, r)
	}
}
