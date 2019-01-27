package handlers

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

// GetNonceFunc should return an effective nonce for calibration (usually retrieved from an EThereum)
type GetNonceFunc func(ctx context.Context, chainID *big.Int, a common.Address) (uint64, error)

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
		defer func() {
			err := nm.Unlock(chainID, a, lockSig)
			if err != nil {
				ctx.Error(err)
			}
		}()

		// Retrieve nonce
		nonce, inCache, err := nm.GetNonce(chainID, a)
		if err != nil {
			ctx.AbortWithError(err)
			return
		}

		// If nonce was not in cache, we calibrate it by reading nonce from chain
		if !inCache {
			nonce, err = getChainNonce(context.Background(), chainID, *a)
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
