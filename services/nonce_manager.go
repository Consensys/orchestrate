package services

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// NonceManager allows to safely manipulate a nonce by locking/unlocking it
type NonceManager interface {
	// Get read nonce value (does not acquire lock), it should indicate if nonce was available or not
	GetNonce(chainID *big.Int, a *common.Address) (uint64, bool, error)

	// Set read nonce value (does not acquire lock), it should indicate if nonce was available or not
	SetNonce(chainID *big.Int, a *common.Address, v uint64) error

	// Lock nonce
	Lock(chainID *big.Int, a *common.Address) (string, error)

	// Unlock nonce
	Unlock(chainID *big.Int, a *common.Address) (string, error)
}
