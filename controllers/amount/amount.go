package amount

import (
	"context"
	"math/big"

	"gitlab.com/ConsenSys/client/fr/core-stack/infra/faucet.git/faucet"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/faucet.git/types"
)

// Controller is a Controller that set the Amount to be transfered
type Controller struct {
	conf *Config
}

// NewController creates a new BlackList controller
func NewController(conf *Config) *Controller {
	return &Controller{
		conf: conf,
	}
}

// Control apply BlackList controller on a credit function
func (ctrl *Controller) Control(credit faucet.CreditFunc) faucet.CreditFunc {
	return func(ctx context.Context, r *types.Request) (*big.Int, bool, error) {
		if ctrl.conf.Amount.Text(10) == "0" {
			return big.NewInt(0), false, nil
		}

		return credit(ctx, &types.Request{
			ChainID:     r.ChainID,
			Creditor:    r.Creditor,
			Beneficiary: r.Beneficiary,
			Amount:      ctrl.conf.Amount,
		})
	}
}
