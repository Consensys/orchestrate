package nonceattributor

import (
	"strconv"

	log "github.com/sirupsen/logrus"
	hnonce "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/nonce"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/nonce/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/nonce"
)

// Handler creates and return an handler for nonce
func Nonce(nm nonce.Attributor, ec hnonce.EthClient) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// TODO: Use the job internal data to do this check once envelope is refactored
		if txctx.Envelope.IsChildJob() {
			txctx.Logger = txctx.Logger.WithFields(log.Fields{
				"nonce": txctx.Envelope.GetNonceString(),
			})
			txctx.Logger.Debug("reusing nonce for child job")
			return
		}

		txctx.Logger.WithField("envelope_id", txctx.Envelope.GetID()).Debugf("nonce handler starts")
		if txctx.Envelope.IsOneTimeKeySignature() {
			txctx.Logger = txctx.Logger.WithFields(log.Fields{
				"chainID":      txctx.Envelope.GetChainIDString(),
				"nonce":        0,
				"one-time-key": true,
			})

			// Set nonce
			_ = txctx.Envelope.SetNonce(0)
			return
		}

		if txctx.Envelope.IsEthSendTesseraPrivateTransaction() {
			_ = txctx.Envelope.SetNonce(0)
			return
		}

		// Retrieve chainID and sender address
		sender, err := txctx.Envelope.GetFromAddress()
		if err != nil {
			_ = txctx.AbortWithError(err).ExtendComponent(component)
			return
		}

		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"from":    sender.Hex(),
			"chainID": txctx.Envelope.GetChainIDString(),
		})

		// Nonce to attribute to tx
		var n uint64
		// Compute nonce key for nonce manager processing
		nonceKey := txctx.Envelope.PartitionKey()

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
				var pendingNonce uint64
				if txctx.Envelope.IsEeaSendPrivateTransaction() {
					pendingNonce, err = utils.EEAGetNonce(txctx.Context(), ec, txctx.Envelope, url)
					if err != nil {
						e := txctx.AbortWithError(err).ExtendComponent(component)
						txctx.Logger.WithError(e).Errorf("nonce: could not read private eea nonce from chain")
						return
					}
				} else {
					pendingNonce, err = utils.GetNonce(txctx.Context(), ec, txctx.Envelope, url)
					if err != nil {
						e := txctx.AbortWithError(err).ExtendComponent(component)
						txctx.Logger.WithError(e).Errorf("nonce: could not read nonce from chain")
						return
					}
				}

				n = pendingNonce
			} else {
				n = lastAttributed + 1
			}
		}

		// Set nonce
		_ = txctx.Envelope.SetNonce(n)
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"nonce":     n,
			"nonce_key": nonceKey,
		})
		txctx.Logger.Debug("nonce: nonce calibrated successfully")

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
