package controls

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/containous/traefik/v2/pkg/log"
	ethcommon "github.com/ethereum/go-ethereum/common"
	stripedmutex "github.com/nmvalera/striped-mutex"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/chain-registry"
)

const cooldownComponent = "faucet.control.cooldown"

// Controller that forces a minimum time interval between 2 credits
type CooldownControl struct {
	mux            *stripedmutex.StripedMutex
	lastAuthorized *sync.Map
}

// NewController creates a CoolDown controller
func NewCooldownControl() *CooldownControl {
	return &CooldownControl{
		lastAuthorized: &sync.Map{},
		mux:            stripedmutex.New(100),
	}
}

// Control apply CoolDown controller on a credit function
func (ctrl *CooldownControl) Control(ctx context.Context, req *types.Request) error {
	if len(req.Candidates) == 0 {
		return errors.FaucetWarning("no faucet candidates").ExtendComponent(cooldownComponent)
	}

	// If still cooling down we invalid credit
	for key, candidate := range req.Candidates {
		ctrl.lock(key, req.Beneficiary)
		defer ctrl.unlock(key, req.Beneficiary)

		if ctrl.IsCoolingDown(key, req.Beneficiary, candidate.Cooldown) {
			log.FromContext(ctx).WithField("beneficiary", req.Beneficiary).
				WithField("faucet", key).Debug("candidate removed due to CooldownControl")
			delete(req.Candidates, key)
		}
	}
	if len(req.Candidates) == 0 {
		return errors.FaucetWarning("faucet cooling down").ExtendComponent(cooldownComponent)
	}

	return nil
}

func (ctrl *CooldownControl) OnSelectedCandidate(_ context.Context, faucet *types.Faucet, beneficiary ethcommon.Address) error {
	ctrl.lastAuthorized.Store(mapKey(faucet.UUID, beneficiary), time.Now())
	return nil
}

// IsCoolingDown indicates if faucet is cooling down
func (ctrl *CooldownControl) IsCoolingDown(faucetID string, beneficiary ethcommon.Address, delay time.Duration) bool {
	lastAuthorized, _ := ctrl.lastAuthorized.LoadOrStore(mapKey(faucetID, beneficiary), time.Time{})
	return time.Since(lastAuthorized.(time.Time)) < delay
}

func (ctrl *CooldownControl) lock(faucetID string, a ethcommon.Address) {
	ctrl.mux.Lock(mapKey(faucetID, a))
}

func (ctrl *CooldownControl) unlock(faucetID string, a ethcommon.Address) {
	ctrl.mux.Unlock(mapKey(faucetID, a))
}

func mapKey(faucetID string, a ethcommon.Address) string {
	return fmt.Sprintf("%s@%s", a.Hex(), faucetID)
}
