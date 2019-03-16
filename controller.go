package faucet

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/striped-mutex.git"
)

// BlackList is a controller that invalid black listed address
type BlackList map[string]struct{}

// NewBlackList creates a new BlackList controller
func NewBlackList(chains []*big.Int, addresses []common.Address) BlackList {
	if len(chains) != len(addresses) {
		panic("BlackList can not be initialized")
	}
	var bl BlackList = make(map[string]struct{})
	for i, chainID := range chains {
		bl[computeKey(chainID, addresses[i])] = struct{}{}
	}
	return bl
}

// IsBlackListed indicates if a user is black listed
func (bl BlackList) IsBlackListed(key string) bool {
	_, ok := bl[key]
	return ok
}

// Control apply BlackList controller on a credit function
func (bl *BlackList) Control(f CreditFunc) CreditFunc {
	return func(ctx context.Context, r *services.FaucetRequest) (*big.Int, bool, error) {
		key := computeKey(r.ChainID, r.Address)
		if bl.IsBlackListed(key) {
			return big.NewInt(0), false, nil
		}
		return f(ctx, r)
	}
}

// CoolDown forces a minimum time interval between 2 credits
type CoolDown struct {
	delay time.Duration

	lastAuthorized *sync.Map
	mux            *stripedmutex.StripedMutex
}

// NewCoolDown creates a CoolDown controller
func NewCoolDown(delay time.Duration, stripes uint) *CoolDown {
	return &CoolDown{
		delay:          delay,
		lastAuthorized: &sync.Map{},
		mux:            stripedmutex.New(stripes),
	}
}

// IsCoolingDown indicates if faucet is cooling doan
func (cd *CoolDown) IsCoolingDown(key string) bool {
	lastAuthorized, _ := cd.lastAuthorized.LoadOrStore(key, time.Time{})
	if time.Now().Sub(lastAuthorized.(time.Time)) < cd.delay {
		return true
	}
	return false
}

// Control apply CoolDosn controller on a credit function
func (cd *CoolDown) Control(f CreditFunc) CreditFunc {
	return func(ctx context.Context, r *services.FaucetRequest) (*big.Int, bool, error) {
		key := computeKey(r.ChainID, r.Address)
		cd.mux.Lock(key)
		defer cd.mux.Unlock(key)

		// If still cooling down we invalid credit
		if cd.IsCoolingDown(key) {
			return big.NewInt(0), false, nil
		}

		// Credit
		amount, ok, err := f(ctx, r)

		// If credit properly occured we update lastAuthorized date
		if ok {
			cd.lastAuthorized.Store(key, time.Now())
		}
		return amount, ok, err
	}

}

// BalanceAtFunc should return a balance
type BalanceAtFunc func(ctx context.Context, chainID *big.Int, a common.Address) (*big.Int, error)

// MaxBalance is a controller that ensures an address can not be credit above a given limit
type MaxBalance struct {
	max       *big.Int
	balanceAt BalanceAtFunc
}

// NewMaxBalance creates a new max balance controller
func NewMaxBalance(max *big.Int, balanceAt BalanceAtFunc) *MaxBalance {
	return &MaxBalance{
		max:       max,
		balanceAt: balanceAt,
	}
}

// Control apply MaxBalance controller on a credit function
func (mb *MaxBalance) Control(f CreditFunc) CreditFunc {
	return func(ctx context.Context, r *services.FaucetRequest) (*big.Int, bool, error) {
		balance, err := mb.balanceAt(ctx, r.ChainID, r.Address)

		if err != nil {
			return big.NewInt(0), false, err
		}

		if balance.Add(balance, r.Value).Cmp(mb.max) >= 0 {
			// Do not credit if final balance would be greater than max authorized
			return big.NewInt(0), false, nil
		}

		// Credit is valid
		return f(ctx, r)
	}
}
