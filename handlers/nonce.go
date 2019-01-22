package handlers

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

// GetNonceFunc is a function which given an eth address and a chain ID returns a nonce
type GetNonceFunc func(chainID *big.Int, a *common.Address) (uint64, error)

type ethClient interface {
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
}

// GetChainNonce returns a function to get nonce initial values from an Eth client
func GetChainNonce(ec ethClient) GetNonceFunc {
	return func(chainID *big.Int, a *common.Address) (uint64, error) {
		return ec.PendingNonceAt(context.Background(), *a)
	}
}

// NonceHandler creates and return an handler for nonce
func NonceHandler(nm services.NonceManager, getChainNonce GetNonceFunc) types.HandlerFunc {
	return func(ctx *types.Context) {
		// Retrieve chainID and sender address
		chainID, a := ctx.T.Chain().ID, ctx.T.Sender().Address

		// Get the lock for chainID and sender address
		lockSig, err := nm.Lock(chainID, a)
		if err != nil {
			ctx.AbortWithError(err)
			return
		}
		defer nm.Unlock(chainID, a, lockSig)

		// Get the nonce from cache
		nonce, inCache, err := nm.GetNonce(chainID, a)
		if err != nil {
			ctx.AbortWithError(err)
			return
		}

		// If the nonce was not in the cache, get it from chain
		if inCache == false {
			nonce, err = getChainNonce(chainID, a)
			if err != nil {
				ctx.AbortWithError(err)
				return
			}
		}

		// Set Nonce value on Trace
		ctx.T.Tx().SetNonce(nonce)

		// Execute pending handlers (note that we do not release lock while executing pending handlers)
		ctx.Next()

		// Increment nonce in Manager
		// TODO: we should ensure pending handlers have correctly executed before incrementing
		err = nm.SetNonce(chainID, a, nonce+1)
		if err != nil {
			// TODO: handle error
			ctx.AbortWithError(err)
			return
		}
	}
}
