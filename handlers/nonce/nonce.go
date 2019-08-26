package nonce

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/nonce.git/nonce"
)

// GetNonceFunc should return an effective nonce for calibration (usually retrieved from an EThereum)
type GetNonceFunc func(ctx context.Context, chainID *big.Int, a common.Address) (uint64, error)

// Handler creates and return an handler for nonce
func Handler(nc nonce.Nonce, getChainNonce GetNonceFunc) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// Retrieve chainID and sender address
		chainID := txctx.Envelope.GetChain().ID()

		a := txctx.Envelope.GetFrom().Address()

		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"tx.sender":   a.Hex(),
			"chain.id":    chainID.String(),
			"metadata.id": txctx.Envelope.GetMetadata().GetId(),
		})

		// Get the lock for chainID and sender address
		lockSig, err := nc.Lock(chainID, &a)
		if err != nil {
			e := txctx.AbortWithError(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("nonce: could not acquire nonce lock")
			return
		}
		defer func() {
			err = nc.Unlock(chainID, &a, lockSig)
			if err != nil {
				e := txctx.AbortWithError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Warnf("nonce: could not release nonce lock")
				return
			}
		}()

		// Retrieve nonce
		n, inCache, err := nc.Get(chainID, &a)
		if err != nil {
			e := txctx.AbortWithError(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("nonce: could not get nonce from cache")
			return
		}

		// If nonce was not in cache, we calibrate it by reading nonce from chain
		if !inCache {
			txctx.Logger.Debugf("nonce: not in cache, get from chain")
			n, err = getChainNonce(txctx.Context(), chainID, a)
			if err != nil {
				e := txctx.AbortWithError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Errorf("nonce: could not get nonce from chain")
				return
			}
		}

		// Set Nonce value on Envelope
		if txctx.Envelope.GetTx() == nil {
			txctx.Envelope.Tx = &ethereum.Transaction{TxData: &ethereum.TxData{}}
		} else if txctx.Envelope.GetTx().GetTxData() == nil {
			txctx.Envelope.Tx.TxData = &ethereum.TxData{}
		}
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
			e := txctx.AbortWithError(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("nonce: could not set nonce on cache")
			return
		}
	}
}
