package nonce

import (
	"context"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/nonce.git/nonce"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// GetNonceFunc should return an effective nonce for calibration (usually retrieved from an EThereum)
type GetNonceFunc func(ctx context.Context, chainID *big.Int, a common.Address) (uint64, error)

// NonceHandler creates and return an handler for nonce
func NonceHandler(nm nonce.Nonce, getChainNonce GetNonceFunc) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// Retrieve chainID and sender address
		chainID := txctx.Envelope.GetChain().ID()

		a, err := txctx.Envelope.GetSender().Address()

		if err != nil {
			txctx.AbortWithError(err)
			txctx.Logger.WithError(err).Errorf("nonce: could not acquire address from sender")
			return
		}

		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"tx.sender": a.Hex(),
			"chain.id":  chainID,
		})

		// Get the lock for chainID and sender address
		lockSig, err := nm.Lock(chainID, &a)
		if err != nil {
			txctx.AbortWithError(err)
			txctx.Logger.WithError(err).Errorf("nonce: could not acquire nonce lock")
			return
		}
		defer func() {
			err := nm.Unlock(chainID, &a, lockSig)
			if err != nil {
				txctx.Error(err)
				txctx.Logger.WithError(err).Errorf("nonce: could not release nonce lock")
			}
		}()

		// Retrieve nonce
		nonce, idleTime, err := nm.Get(chainID, &a)
		if err != nil {
			txctx.AbortWithError(err)
			txctx.Logger.WithError(err).Errorf("nonce: could not get nonce from cache")
			return
		}

		// If nonce was not in cache, we calibrate it by reading nonce from chain
		if idleTime == -1 {
			txctx.Logger.Debugf("nonce: not in cache, get from chain")
			nonce, err = getChainNonce(txctx.Context(), chainID, a)
			if err != nil {
				txctx.AbortWithError(err)
				txctx.Logger.WithError(err).Errorf("nonce: could not get nonce from chain")
				return
			}
		}

		// If nonce is too old, we calibrate it by reading nonce from chain
		if idleTime > viper.GetInt("redis.nonce.expiration.time") {
			txctx.Logger.Debugf("nonce: cache too old, get from chain")
			nonce, err = getChainNonce(txctx.Context(), chainID, a)
			if err != nil {
				txctx.AbortWithError(err)
				txctx.Logger.WithError(err).Errorf("nonce: could not get nonce from chain")
				return
			}
		}

		// Set Nonce value on Trace
		txctx.Envelope.GetTx().GetTxData().SetNonce(nonce)
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"tx.nonce": nonce,
		})
		txctx.Logger.Debugf("nonce: nonce set")

		// Execute pending handlers (note that we do not release lock while executing pending handlers)
		txctx.Next()

		// Increment nonce in Manager
		// TODO: we should ensure pending handlers have correctly executed before incrementing
		err = nm.Set(chainID, &a, nonce+1)
		if err != nil {
			// TODO: handle error
			txctx.AbortWithError(err)
			txctx.Logger.WithError(err).Errorf("nonce: could not set nonce on cache")
			return
		}
	}
}
