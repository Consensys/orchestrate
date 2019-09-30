package nonce

import (
	"strconv"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/types/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/nonce"
)

// Handler creates and return an handler for nonce
func Nonce(nm nonce.Attributor, ec ethclient.ChainStateReader) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// Retrieve chainID and sender address
		chainID, sender := txctx.Envelope.GetChain().ID(), txctx.Envelope.GetFrom().Address()
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"tx.sender":   sender.Hex(),
			"chain.id":    chainID.String(),
			"metadata.id": txctx.Envelope.GetMetadata().GetId(),
		})

		// Nonce to attribute to tx
		var n uint64
		// Compute nonce key for nonce manager procesing
		nonceKey := string(txctx.In.Key())

		// First check if signal for recovering nonce
		if v, ok := txctx.Envelope.GetMetadataValue("nonce.recovering.expected"); ok {
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

			// Reset signal field
			delete(txctx.Envelope.GetMetadata().GetExtra(), "nonce.recovering.expected")
		} else {
			// No signal for nonce recovery
			// Retrieve last attributed nonce from nonce manager
			lastAttributed, ok, err := nm.GetLastAttributed(nonceKey)
			if err != nil {
				e := txctx.AbortWithError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Errorf("nonce: could not retrieve last attributed nonce")
				return
			}

			// If no nonce is available in nonce manager
			// we calibrate by querying chain
			if !ok {
				txctx.Logger.Debugf("nonce: not in cache, get from chain")

				// Retrieve nonce from chain
				pendingNonce, err := ec.PendingNonceAt(txctx.Context(), chainID, sender)
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
		setNonce(txctx.Envelope, n)
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
				txctx.Logger.WithError(e).Errorf("nonce: could not set nonce on cache")
			}
		}
	}
}

func setNonce(e *envelope.Envelope, n uint64) {
	// Initialiaze Transaction on envelope if needed
	if e.GetTx() == nil {
		e.Tx = &ethereum.Transaction{TxData: &ethereum.TxData{}}
	} else if e.GetTx().GetTxData() == nil {
		e.Tx.TxData = &ethereum.TxData{}
	}
	// Set transaction nonce
	e.GetTx().GetTxData().SetNonce(n)
}
