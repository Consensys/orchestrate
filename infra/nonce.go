package infra

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/striped-mutex.git"
)

// SafeNonce allow to manipulate nonce in a concurrently safe manner
type SafeNonce struct {
	mux   *sync.Mutex
	value uint64
}

// Lock acquire lock
func (n *SafeNonce) Lock() error {
	n.mux.Lock()
	return nil
}

// Unlock release lock
func (n *SafeNonce) Unlock() error {
	n.mux.Unlock()
	return nil
}

// Get retrieve nonce value
// Warning: it does not acquire the lock
func (n *SafeNonce) Get() (uint64, error) {
	return n.value, nil
}

// Set sets nonce value
// Warning: it does not acquire the lock
func (n *SafeNonce) Set(v uint64) error {
	n.value = v
	return nil
}

// CalibrateNonceFunc allows to calibrate managed nonce
type CalibrateNonceFunc func(chainID *big.Int, a common.Address) (uint64, error)

// CacheNonceManager allows to manage nonce
type CacheNonceManager struct {
	mux    *stripedmutex.StripedMutex
	nonces *sync.Map

	calibrate CalibrateNonceFunc
}

// NewCacheNonceManager creates a new cache nonce
func NewCacheNonceManager(calibrate CalibrateNonceFunc, stripes uint) *CacheNonceManager {
	return &CacheNonceManager{
		mux:       stripedmutex.New(stripes),
		nonces:    &sync.Map{},
		calibrate: calibrate,
	}
}

func computeKey(chainID *big.Int, a common.Address) string {
	return fmt.Sprintf("%v-%v", chainID.Text(16), a.Hex())
}

// Obtain return a locked SafeNonce for given chain and address
func (c *CacheNonceManager) Obtain(chainID *big.Int, a common.Address) (services.NonceLocker, error) {
	key := computeKey(chainID, a)
	mux, err := c.mux.GetLock(key)
	if err != nil {
		return nil, err
	}
	// Lock key
	mux.Lock()
	defer mux.Unlock()

	// Retrieve nonce from cache
	n, ok := c.nonces.LoadOrStore(key, &SafeNonce{mux: mux, value: 0})
	rv := n.(*SafeNonce)
	if !ok {
		// If nonce has just been created we compute its initial value
		rv.value, err = c.calibrate(chainID, a)
		if err != nil {
			return rv, err
		}
	}
	return rv, nil
}

// NewEthClientNonceCalibrate returns a function to get nonce initial values from an Eth client
func NewEthClientNonceCalibrate(ec *EthClient) CalibrateNonceFunc {
	return func(chainID *big.Int, a common.Address) (uint64, error) {
		return ec.PendingNonceAt(context.Background(), a)
	}
}
