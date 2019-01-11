package infra

import (
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

// StripedNonceManager allows to manage nonce based on a striped mutex
type StripedNonceManager struct {
	mux    *stripedmutex.StripedMutex
	nonces *sync.Map
}

// NewStripedNonceManager creates a new StripedNonceManager
func NewStripedNonceManager(stripes uint) *StripedNonceManager {
	return &StripedNonceManager{
		mux:    stripedmutex.New(stripes),
		nonces: &sync.Map{},
	}
}

func computeKey(chainID *big.Int, a common.Address) string {
	return fmt.Sprintf("%v-%v", chainID.Text(16), a.Hex())
}

// Obtain return a locked SafeNonce for given chain and address
func (c *StripedNonceManager) Obtain(chainID *big.Int, a common.Address) (services.NonceLocker, bool, error) {
	key := computeKey(chainID, a)
	mux, err := c.mux.GetLock(key)
	if err != nil {
		return nil, false, err
	}
	// Lock key
	mux.Lock()
	defer mux.Unlock()

	// Retrieve nonce from cache
	n, ok := c.nonces.LoadOrStore(key, &SafeNonce{mux: mux, value: 0})
	rv := n.(*SafeNonce)

	return rv, ok, nil
}
