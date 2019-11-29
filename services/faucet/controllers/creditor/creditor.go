package creditor

import (
	"context"
	"math/big"
	"sync"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/faucet"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/types"
)

// Controller is a controller that holds a list of account that should not be credited
type Controller struct {
	creditors *sync.Map
}

// NewController creates a new BlackList controller
func NewController() *Controller {
	return &Controller{
		creditors: &sync.Map{},
	}
}

// SetCreditor set creditor address for a given chain
func (ctrl *Controller) SetCreditor(chainID *big.Int, address ethcommon.Address) {
	ctrl.creditors.Store(chainID.Text(16), address)
}

// Creditor returns creditor for a given chain
func (ctrl *Controller) Creditor(chainID *big.Int) (ethcommon.Address, bool) {
	addr, ok := ctrl.creditors.Load(chainID.Text(16))
	if !ok {
		return ethcommon.Address{}, false
	}
	return addr.(ethcommon.Address), ok
}

// Control apply BlackList controller on a credit function
func (ctrl *Controller) Control(credit faucet.CreditFunc) faucet.CreditFunc {
	return func(ctx context.Context, r *types.Request) (*big.Int, error) {
		creditor, ok := ctrl.Creditor(r.ChainID)
		if !ok {
			// No creditor for the given chain
			return big.NewInt(0), errors.FaucetNotConfiguredWarning("no creditor for the given chain")
		}

		if creditor.Hex() == r.Beneficiary.Hex() {
			// Creditor does auto-credit
			return big.NewInt(0), errors.FaucetSelfCreditWarning("attempt to credit the creditor")
		}

		return credit(ctx, &types.Request{
			ChainID:     r.ChainID,
			Creditor:    creditor,
			Beneficiary: r.Beneficiary,
			Amount:      r.Amount,
		})
	}
}
