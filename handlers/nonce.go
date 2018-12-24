package handlers

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
)

// LockedNonce is an interface for a nonce object that can be manipulated safely
type LockedNonce interface {
	Unlock() error
	Get() (uint64, error)
	Set(v uint64) error
}

// NonceManager is an interface for managing nonces
type NonceManager interface {
	// Lock is expected to lock nonce modifications
	Lock(chainID string, a common.Address) (LockedNonce, error)
}

// SafeNonce provides a value an a lock to manipulate it
type SafeNonce struct {
	value uint64
	mux   *sync.Mutex
}

// Lock acquire lock on SafeNonce
func (n *SafeNonce) Lock() error {
	n.mux.Lock()
	return nil
}

// Unlock release lock on SafeNonce
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

// NewNonceFunc is a function to create new nonces
type NewNonceFunc func(a common.Address) (*SafeNonce, error)

// CacheNonce allows to store mutiple SafeNonce
type CacheNonce struct {
	mux    *sync.RWMutex
	nonces map[string]*SafeNonce

	new NewNonceFunc
}

// Get returns a Nonce from cache (it creates it if not already existing)
func (c *CacheNonce) Get(a common.Address) (*SafeNonce, error) {
	// Acquire Read lock
	c.mux.RLock()

	n, ok := c.nonces[a.Hex()]
	if ok {
		// There is already an entry for pair Address
		c.mux.RUnlock()
		return n, nil
	}
	// There was no entry for address we need to write it
	c.mux.RUnlock()

	// Acquire Write lock
	c.mux.Lock()

	// We ensure no entry has been written while acquiring lock
	n, ok = c.nonces[a.Hex()]
	if ok {
		// Nothing has been written so we can proceed
		c.mux.Unlock()
		return n, nil
	}

	n, err := c.newNonce(a)
	if err != nil {
		return nil, err
	}
	c.mux.Unlock()
	return n, nil
}

func (c *CacheNonce) newNonce(a common.Address) (*SafeNonce, error) {
	n, err := c.new(a)
	if err != nil {
		return nil, err
	}
	c.nonces[a.Hex()] = n
	return n, nil
}

// CacheNonceManager is a NonceManager that relies on internal cache
// Can only handle one chain
type CacheNonceManager struct {
	c *CacheNonce
}

// NewCacheNonceManager creates a new manager
func NewCacheNonceManager(new NewNonceFunc) *CacheNonceManager {
	return &CacheNonceManager{
		c: &CacheNonce{
			&sync.RWMutex{},
			make(map[string]*SafeNonce),
			new,
		},
	}
}

// Lock return a SafeNonce locked
func (m *CacheNonceManager) Lock(chainID string, a common.Address) (LockedNonce, error) {
	n, err := m.c.Get(a)
	if err != nil {
		return nil, err
	}
	n.Lock()
	return n, nil
}

// NonceHandler creates and return an handler for nonce
func NonceHandler(m NonceManager) infra.HandlerFunc {
	return func(ctx *infra.Context) {
		// Retrieve chainId and sender address
		chainID, a := ctx.T.Chain().ID, *ctx.T.Sender().Address

		// Lock nonce
		n, err := m.Lock(chainID, a)
		if err != nil {
			// Deal with nonce error
			return
		}
		defer n.Unlock()

		// Set nonce
		v, _ := n.Get()
		ctx.T.Tx().SetNonce(v)

		// Execute pending handlers
		ctx.Next()

		// Increment nonce
		// TODO: we should ensure pending handlers have correctly executed before incrementing
		n.Set(ctx.T.Tx().Nonce() + 1)
	}
}
