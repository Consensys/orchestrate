package cooldown

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"

	ethcommon "github.com/ethereum/go-ethereum/common"
	stripedmutex "github.com/nmvalera/striped-mutex"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/faucet"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/types"
)

// Controller that forces a minimum time interval between 2 credits
type Controller struct {
	mux            *stripedmutex.StripedMutex
	lastAuthorized *sync.Map
}

// NewController creates a CoolDown controller
func NewController() *Controller {
	return &Controller{
		lastAuthorized: &sync.Map{},
		mux:            stripedmutex.New(100),
	}
}

func mapKey(faucetID string, a ethcommon.Address) string {
	return fmt.Sprintf("%s@%s", a.Hex(), faucetID)
}

// IsCoolingDown indicates if faucet is cooling down
func (ctrl *Controller) IsCoolingDown(faucetID string, a ethcommon.Address, delay time.Duration) bool {
	lastAuthorized, _ := ctrl.lastAuthorized.LoadOrStore(mapKey(faucetID, a), time.Time{})
	return time.Since(lastAuthorized.(time.Time)) < delay
}

func (ctrl *Controller) lock(faucetID string, a ethcommon.Address) {
	ctrl.mux.Lock(mapKey(faucetID, a))
}

func (ctrl *Controller) unlock(faucetID string, a ethcommon.Address) {
	ctrl.mux.Unlock(mapKey(faucetID, a))
}

// Authorized is called to indicate that a credit has been authorized
func (ctrl *Controller) Authorized(faucetID string, a ethcommon.Address) {
	ctrl.lastAuthorized.Store(mapKey(faucetID, a), time.Now())
}

// Control apply CoolDown controller on a credit function
func (ctrl *Controller) Control(credit faucet.CreditFunc) faucet.CreditFunc {
	return func(ctx context.Context, r *types.Request) (*big.Int, error) {
		if len(r.FaucetsCandidates) == 0 {
			return nil, errors.FaucetWarning("no faucet candidates").ExtendComponent(component)
		}

		// If still cooling down we invalid credit
		for key, candidate := range r.FaucetsCandidates {
			ctrl.lock(key, r.Beneficiary)
			defer ctrl.unlock(key, r.Beneficiary)

			if ctrl.IsCoolingDown(key, r.Beneficiary, candidate.Cooldown) {
				delete(r.FaucetsCandidates, key)
			}
		}
		if len(r.FaucetsCandidates) == 0 {
			return nil, errors.FaucetWarning("faucet cooling down").ExtendComponent(component)
		}

		// Credit
		amount, err := credit(ctx, r)

		// If credit occurred we update
		if err == nil {
			ctrl.Authorized(r.ElectedFaucet, r.Beneficiary)
		}

		return amount, err
	}
}
