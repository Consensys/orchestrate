package handlers

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
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

		ctx.Logger = ctx.Logger.WithFields(log.Fields{
			"tx.sender": a.Hex(),
			"chain.id":  chainID.Text(16),
		})

		// Get the lock for chainID and sender address
		lockSig, err := nm.Lock(chainID, a)
		if err != nil {
			ctx.AbortWithError(err)
			ctx.Logger.WithError(err).Errorf("nonce: could not acquire nonce lock")
			return
		}
		defer func() {
			err := nm.Unlock(chainID, a, lockSig)
			if err != nil {
				ctx.Error(err)
				ctx.Logger.WithError(err).Errorf("nonce: could not release nonce lock")
			}
		}()

		// Retrieve nonce
		nonce, idleTime, err := nm.GetNonce(chainID, a)
		if err != nil {
			ctx.AbortWithError(err)
			ctx.Logger.WithError(err).Errorf("nonce: could not get nonce from cache")
			return
		}

		// If nonce was not in cache, we calibrate it by reading nonce from chain
		if idleTime == -1 {
			ctx.Logger.Debugf("nonce: not in cache, get from chain")
			nonce, err = getChainNonce(context.Background(), chainID, *a)
			if err != nil {
				ctx.AbortWithError(err)
				ctx.Logger.WithError(err).Errorf("nonce: could not get nonce from chain")
				return
			}
		}

		// If nonce is too old, we calibrate it by reading nonce from chain
		if idleTime > 3 {
			ctx.Logger.Debugf("nonce: cache too old, get from chain")
			nonce, err = getChainNonce(context.Background(), chainID, *a)
			if err != nil {
				ctx.AbortWithError(err)
				ctx.Logger.WithError(err).Errorf("nonce: could not get nonce from chain")
				return
			}
		}

		// Set Nonce value on Trace
		ctx.T.Tx().SetNonce(nonce)
		ctx.Logger = ctx.Logger.WithFields(log.Fields{
			"tx.nonce": nonce,
		})
		ctx.Logger.Debugf("nonce: nonce set")

		// Execute pending handlers (note that we do not release lock while executing pending handlers)
		ctx.Next()

		// Increment nonce in Manager
		// TODO: we should ensure pending handlers have correctly executed before incrementing
		err = nm.SetNonce(chainID, a, nonce+1)
		if err != nil {
			// TODO: handle error
			ctx.AbortWithError(err)
			ctx.Logger.WithError(err).Errorf("nonce: could not set nonce on cache")
			return
		}
	}
}
