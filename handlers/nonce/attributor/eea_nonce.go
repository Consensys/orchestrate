package nonceattributor

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/nonce"
)

// Handler creates and return an handler for marking tx nonce (Only for Orion private transactions)
// @TODO Remove once Orion tx is spit into two jobs
//  https://app.zenhub.com/workspaces/orchestrate-5ea70772b186e10067f57842/issues/pegasyseng/orchestrate/253
func EEANonce(nm nonce.Attributor, ec ethclient.ChainStateReader) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Logger.WithField("envelope_id", txctx.Envelope.GetID()).Debugf("eea_nonce handler starts")
		if !txctx.Envelope.IsEeaSendPrivateTransaction() {
			return
		}

		if txctx.Envelope.IsOneTimeKeySignature() {
			_ = txctx.Envelope.SetEEAMarkingTxNonce(0)
			txctx.Logger = txctx.Logger.WithFields(log.Fields{
				"eea_nonce": 0,
			})
			return
		}

		// Nonce to attribute to tx
		var n uint64
		// Compute nonce key for nonce manager processing
		nonceKey := calcNonceKey(txctx.Envelope)

		if v := txctx.Envelope.GetInternalLabelsValue("nonce.recovering.expected"); v != "" {
			url, err := proxy.GetURL(txctx)
			if err != nil {
				return
			}

			expectedNonce, err := ec.PendingNonceAt(txctx.Context(), url, txctx.Envelope.MustGetFromAddress())
			if err != nil {
				e := txctx.AbortWithError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Errorf("eea_nonce: calibrating nonce from chain after recovering")
				return
			}

			n = expectedNonce
		} else {
			// No signal for nonce recovery
			// Retrieve last attributed nonce from nonce manager
			lastAttributed, ok, err := nm.GetLastAttributed(nonceKey)
			if err != nil {
				e := txctx.AbortWithError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Errorf("eea_nonce: could not load last attributed nonce")
				return
			}

			// If no nonce is available in nonce manager
			// we calibrate by querying chain
			if !ok {
				txctx.Logger.Debugf("eea_nonce: calibrating nonce from chain")
				url, err := proxy.GetURL(txctx)
				if err != nil {
					return
				}

				// Retrieve nonce from chain
				pendingNonce, err := ec.PendingNonceAt(txctx.Context(), url, txctx.Envelope.MustGetFromAddress())
				if err != nil {
					e := txctx.AbortWithError(err).ExtendComponent(component)
					txctx.Logger.WithError(e).Errorf("eea_nonce: could not read nonce from chain")
					return
				}

				n = pendingNonce
			} else {
				n = lastAttributed + 1
			}
		}

		// Set nonce
		_ = txctx.Envelope.SetEEAMarkingTxNonce(n)
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"eea_nonce":     n,
			"eea_nonce_key": nonceKey,
		})

		// Execute pending handlers
		txctx.Next()

		// If pending handlers executed correctly we increment nonce
		if len(txctx.Envelope.GetErrors()) == 0 {
			// Increment nonce
			err := nm.SetLastAttributed(nonceKey, n)
			if err != nil {
				e := errors.FromError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Errorf("nonce: could not store last attributed nonce")
			}
		}
	}
}

func calcNonceKey(e *tx.Envelope) string {
	var chainKey string
	if e.GetChainID() != nil {
		chainKey = e.GetChainID().String()
	} else {
		chainKey = e.GetChainName()
	}

	return fmt.Sprintf("%v@%v", e.GetFromString(), chainKey)
}
