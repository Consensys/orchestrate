package maxbalance

import (
	"context"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/faucet.git/faucet"
)

// BalanceAtFunc should return a balance
type BalanceAtFunc func(ctx context.Context, chainID *big.Int, a ethcommon.Address, blocknumber *big.Int) (*big.Int, error)

// Controller is a controller that ensures an address can not be credit above a given limit
type Controller struct {
	conf *Config
}

// NewController creates a new max balance controller
func NewController(conf *Config) *Controller {
	return &Controller{
		conf,
	}
}

// Control apply MaxBalance controller on a credit function
func (ctrl *Controller) Control(credit faucet.CreditFunc) faucet.CreditFunc {
	return func(ctx context.Context, r *faucet.Request) (*big.Int, bool, error) {
		// Retrieve account balance
		balance, err := ctrl.conf.BalanceAt(ctx, r.ChainID, r.Address, nil)
		if err != nil {
			return big.NewInt(0), false, err
		}

		// Ensure MaxBalance is repected
		if balance.Add(balance, r.Value).Cmp(ctrl.conf.MaxBalance) >= 0 {
			// Do not credit if final balance would be superior to max authorized
			return big.NewInt(0), false, nil
		}

		return credit(ctx, r)
	}
}
