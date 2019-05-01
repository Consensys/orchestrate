package mock

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

var(
	lockSig = "mokcLockId"
)

// Nonce is a mock Nonce
type Nonce struct{}

// NewNonce creates a new mock faucet
func NewNonce() *Nonce {
	return &Nonce{}
}

// Get read nonce value (does not acquire lock), it should indicate if nonce was available or not
func (nm *Nonce) Get(chainID *big.Int, a *common.Address) (uint64, int, error) {
	return 0, -1, nil // idleTime == -1, meaning the nonce is not in the cache
}

// Set read nonce value (does not acquire lock), it should indicate if nonce was available or not
func (nm *Nonce) Set(chainID *big.Int, a *common.Address, v uint64) error {
	return nil
}

// Lock nonce
func (nm *Nonce) Lock(chainID *big.Int, a *common.Address) (string, error) {
	return lockSig, nil
}

// Unlock nonce
func (nm *Nonce) Unlock(chainID *big.Int, a *common.Address, lockSig string) error {
	return nil
}