package blacklist

import (
	"context"
	"math/big"
	"sync"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/faucet.git/faucet"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/utils"
)

// Controller is a controller that holds a list of account that should not be credited
type Controller struct {
	blacklist *sync.Map
}

// NewController creates a new BlackList controller
func NewController() *Controller {
	return &Controller{
		blacklist: &sync.Map{},
	}
}

// BlackList  an account on a given chain
func (ctrl *Controller) BlackList(chainID *big.Int, address ethcommon.Address) {
	ctrl.blacklist.Store(utils.ToChainAccountKey(chainID, address), struct{}{})
}

// IsBlackListed indicates if a user is black listed
func (ctrl *Controller) IsBlackListed(chainID *big.Int, address ethcommon.Address) bool {
	key := utils.ToChainAccountKey(chainID, address)
	_, ok := ctrl.blacklist.Load(key)
	return ok
}

// Control apply BlackList controller on a credit function
func (ctrl *Controller) Control(credit faucet.CreditFunc) faucet.CreditFunc {
	return func(ctx context.Context, r *faucet.Request) (*big.Int, bool, error) {
		if ctrl.IsBlackListed(r.ChainID, r.Address) {
			return big.NewInt(0), false, nil
		}
		return credit(ctx, r)
	}
}
