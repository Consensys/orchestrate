package services

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// NonceLocker allows to safely manipulate a nonce by locking/unlocking it
type NonceLocker interface {
	Lock() error
	Get() (uint64, error)
	Set(v uint64) error
	Unlock() error
}

// NonceManager is an interface for fine grain management of nonce by key
type NonceManager interface {
	// Return a NonceLocker
	Obtain(chainID *big.Int, a common.Address) (NonceLocker, bool, error)
}

// NonceCalibrator is a function expected to return a calibrated nonce value
// For example when we meet a pair chainID, address for the first time
type NonceCalibrator func(chainID *big.Int, a common.Address) (uint64, error)
