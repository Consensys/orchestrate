package cooldown

import (
	"context"
	"math/big"
	"sync"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	stripedmutex "github.com/nmvalera/striped-mutex"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/faucet/faucet"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/faucet/types"
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

// IsCoolingDown indicates if faucet is cooling down
func (ctrl *Controller) IsCoolingDown(chainID *big.Int, a ethcommon.Address) bool {
	key := utils.ToChainAccountKey(chainID, a)
	lastAuthorized, _ := ctrl.lastAuthorized.LoadOrStore(key, time.Time{})
	return time.Since(lastAuthorized.(time.Time)) < ctrl.conf.Delay
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
			return big.NewInt(0), false, errors.FaucetWarning("faucet cooling down").ExtendComponent(component)
		}

		// Credit
		amount, ok, err := credit(ctx, r)

		// If credit occurred we update
		if ok {
			ctrl.Authorized(r.ChainID, r.Beneficiary)
		}

		return amount, ok, err
	}
}
