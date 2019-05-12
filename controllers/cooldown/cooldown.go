package cooldown

import (
	"context"
	"math/big"
	"sync"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/faucet.git/faucet"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/faucet.git/types"
	stripedmutex "gitlab.com/ConsenSys/client/fr/core-stack/striped-mutex.git"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/utils"
)

// Controller that forces a minimum time interval between 2 credits
type Controller struct {
	conf *Config

	mux            *stripedmutex.StripedMutex
	lastAuthorized *sync.Map
}

// NewController creates a CoolDown controller
func NewController(conf *Config) *Controller {
	return &Controller{
		conf:           conf,
		lastAuthorized: &sync.Map{},
		mux:            stripedmutex.New(100),
	}
}

// IsCoolingDown indicates if faucet is cooling doan
func (ctrl *Controller) IsCoolingDown(chainID *big.Int, a ethcommon.Address) bool {
	key := utils.ToChainAccountKey(chainID, a)
	lastAuthorized, _ := ctrl.lastAuthorized.LoadOrStore(key, time.Time{})
	if time.Now().Sub(lastAuthorized.(time.Time)) < ctrl.conf.Delay {
		return true
	}
	return false
}

func (ctrl *Controller) lock(chainID *big.Int, a ethcommon.Address) {
	key := utils.ToChainAccountKey(chainID, a)
	ctrl.mux.Lock(key)
}

func (ctrl *Controller) unlock(chainID *big.Int, a ethcommon.Address) {
	key := utils.ToChainAccountKey(chainID, a)
	ctrl.mux.Unlock(key)
}

// Authorized is called to indicate that a credit has been authorized
func (ctrl *Controller) Authorized(chainID *big.Int, a ethcommon.Address) {
	key := utils.ToChainAccountKey(chainID, a)
	ctrl.lastAuthorized.Store(key, time.Now())
}

// Control apply CoolDown controller on a credit function
func (ctrl *Controller) Control(credit faucet.CreditFunc) faucet.CreditFunc {
	return func(ctx context.Context, r *types.Request) (*big.Int, bool, error) {
		ctrl.lock(r.ChainID, r.Beneficiary)
		defer ctrl.unlock(r.ChainID, r.Beneficiary)

		// If still cooling down we invalid credit
		if ctrl.IsCoolingDown(r.ChainID, r.Beneficiary) {
			return big.NewInt(0), false, nil
		}

		// Credit
		amount, ok, err := credit(ctx, r)

		// If credit occured we update
		if ok {
			ctrl.Authorized(r.ChainID, r.Beneficiary)
		}

		return amount, ok, err
	}
}
