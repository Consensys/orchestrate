package nonceattributor

import (
	"strconv"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/nonce"
)

// Handler creates and return an handler for nonce
func Nonce(nm nonce.Attributor, ec ethclient.ChainStateReader) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// Retrieve chainID and sender address
		sender, err := txctx.Envelope.GetFromAddress()
		if err != nil {
			_ = txctx.AbortWithError(err).ExtendComponent(component)
			return
		}

		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"from":    sender.Hex(),
			"chainID": txctx.Envelope.GetChainIDString(),
			"id":      txctx.Envelope.GetID(),
		})

		// Nonce to attribute to tx
		var n uint64
		// Compute nonce key for nonce manager processing
		nonceKey := string(txctx.In.Key())

		// First check if signal for recovering nonce
		if v := txctx.Envelope.GetInternalLabelsValue("nonce.recovering.expected"); v != "" {
			expectedNonce, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				e := txctx.AbortWithError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Errorf("nonce: could not extract metadata")
				return
			}

			txctx.Logger.WithFields(log.Fields{
				"nonce.expected": expectedNonce,
			}).Warnf("nonce: recalibrate nonce following recovery signal")

			n = expectedNonce
		} else {
			// No signal for nonce recovery
			// Retrieve last attributed nonce from nonce manager
			lastAttributed, ok, err := nm.GetLastAttributed(nonceKey)
			if err != nil {
				e := txctx.AbortWithError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Errorf("nonce: could not load last attributed nonce")
				return
			}

			// If no nonce is available in nonce manager
			// we calibrate by querying chain
			if !ok {
				txctx.Logger.Debugf("nonce: calibrating nonce from chain")
				url, err := proxy.GetURL(txctx)
				if err != nil {
					return
				}

				// Retrieve nonce from chain
				pendingNonce, err := ec.PendingNonceAt(txctx.Context(), url, sender)
				if err != nil {
					e := txctx.AbortWithError(err).ExtendComponent(component)
					txctx.Logger.WithError(e).Errorf("nonce: could not read nonce from chain")
					return
				}

				n = pendingNonce
			} else {
				n = lastAttributed + 1
			}
		}

		// Set nonce
		_ = txctx.Envelope.SetNonce(n)
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"tx.nonce": n,
		})

		// Execute pending handlers
		txctx.Next()

		// If pending handlers executed correctly we increment nonce
		if len(txctx.Envelope.GetErrors()) == 0 {
			// Increment nonce
			err := nm.SetLastAttributed(nonceKey, n)
			if err != nil {
				// An error here means that we probably lost connection with NonceManager underlying cache.
				// TODO: A retry strategy should be implemented on the nonce manager to make this scenario rare
				//
				// At this point pending handlers have executed correctly
				// meaning that envelope has already been successfully produced to the next topic
				//
				// So we log the error and process next envelope
				e := errors.FromError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Errorf("nonce: could not store last attributed nonce")
			}
		}
	}
}
