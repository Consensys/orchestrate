package nonce

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/nonce.git/nonce"
)

// GetNonceFunc should return an effective nonce for calibration (usually retrieved from an EThereum)
type GetNonceFunc func(ctx context.Context, chainID *big.Int, a common.Address) (uint64, error)

// Handler creates and return an handler for nonce
func Handler(nc nonce.Nonce, getChainNonce GetNonceFunc) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// Retrieve chainID and sender address
		chainID := txctx.Envelope.GetChain().ID()

		a, err := txctx.Envelope.GetSender().Address()

		if err != nil {
			txctx.Logger.WithError(err).Errorf("nonce: could not acquire address from sender")
			_ = txctx.AbortWithError(err)
			return
		}

		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"tx.sender": a.Hex(),
			"chain.id":  chainID,
		})

		// Get the lock for chainID and sender address
		lockSig, err := nc.Lock(chainID, &a)
		if err != nil {
			txctx.Logger.WithError(err).Errorf("nonce: could not acquire nonce lock")
			_ = txctx.AbortWithError(err)
			return
		}
		defer func() {
			er := nc.Unlock(chainID, &a, lockSig)
			if er != nil {
				txctx.Logger.WithError(err).Errorf("nonce: could not release nonce lock")
				_ = txctx.Error(er)
			}
		}()

		// Retrieve nonce
		n, idleTime, err := nc.Get(chainID, &a)
		if err != nil {
			txctx.Logger.WithError(err).Errorf("nonce: could not get nonce from cache")
			_ = txctx.AbortWithError(err)
			return
		}

		// If nonce was not in cache, we calibrate it by reading nonce from chain
		if idleTime == -1 {
			txctx.Logger.Debugf("nonce: not in cache, get from chain")
			n, err = getChainNonce(txctx.Context(), chainID, a)
			if err != nil {
				txctx.Logger.WithError(err).Errorf("nonce: could not get nonce from chain")
				_ = txctx.AbortWithError(err)
				return
			}
		}

		// If nonce is too old, we calibrate it by reading nonce from chain
		if idleTime > viper.GetInt("redis.nonce.expiration.time") {
			txctx.Logger.Debugf("nonce: cache too old, get from chain")
			n, err = getChainNonce(txctx.Context(), chainID, a)
			if err != nil {
				txctx.Logger.WithError(err).Errorf("nonce: could not get nonce from chain")
				_ = txctx.AbortWithError(err)
				return
			}
		}

		// Set Nonce value on Trace
		txctx.Envelope.GetTx().GetTxData().SetNonce(n)
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"tx.nonce": n,
		})
		txctx.Logger.Debugf("nonce: nonce set")

		// Execute pending handlers (note that we do not release lock while executing pending handlers)
		txctx.Next()

		// Increment nonce in Manager
		// TODO: we should ensure pending handlers have correctly executed before incrementing
		err = nc.Set(chainID, &a, n+1)
		if err != nil {
			// TODO: handle error
			txctx.Logger.WithError(err).Errorf("nonce: could not set nonce on cache")
			_ = txctx.AbortWithError(err)
			return
		}
	}
}
