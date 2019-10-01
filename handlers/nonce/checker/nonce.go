package noncechecker

import (
	"strconv"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/types/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/nonce"
)

// Checker creates an handler responsible to check transaction nonce value
// It makes sure that transactions with invalid nonce are not processed
func Checker(nm nonce.Sender, ec ethclient.ChainStateReader) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// Retrieve chainID and sender address
		chainID, sender := txctx.Envelope.GetChain().ID(), txctx.Envelope.GetFrom().Address()
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"tx.sender":   sender.Hex(),
			"chain.id":    chainID.String(),
			"metadata.id": txctx.Envelope.GetMetadata().GetId(),
		})

		// Retrieves nonce key for nonce manager procesing
		nonceKey := string(txctx.In.Key())
		var expectedNonce uint64

		// Retrieve last sent nonce from nonce manager
		lastSent, ok, err := nm.GetLastSent(nonceKey)
		if err != nil {
			e := txctx.AbortWithError(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("nonce: could not retrieve last sent nonce")
			return
		}

		// If no nonce is available (meaning that envelope being processed is the first one for a given sender on a given chain)
		// then we calibrate nonce manager by retrieving nonce from chain
		if ok {
			expectedNonce = lastSent + 1
		} else {
			// Retrieve nonce from chain
			pendingNonce, err := ec.PendingNonceAt(txctx.Context(), chainID, sender)
			if err != nil {
				e := txctx.AbortWithError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Errorf("nonce: could not read nonce from chain")
				return
			}
			txctx.Logger.WithFields(log.Fields{
				"nonce.pending": pendingNonce,
			}).Debugf("nonce: retrieving nonce from chain")

			expectedNonce = pendingNonce
		}

		// Compare expected nonce to attributed nonced
		// to make sure we do not sent a transaction with invalid nonce
		n := txctx.Envelope.GetTx().GetTxData().GetNonce()
		if n != expectedNonce {

			// Transaction has an invalid nonce
			txctx.Logger.WithFields(log.Fields{
				"nonce.expected": expectedNonce,
				"nonce.got":      n,
			}).Warnf("nonce: invalid nonce")

			// Reset tx nonce, hash and raw
			resetTx(txctx.Envelope.GetTx())

			if n > expectedNonce {
				// Nonce is too high

				// Retrieve recovery status
				recovering, err := nm.IsRecovering(nonceKey)
				if err != nil {
					e := txctx.AbortWithError(err).ExtendComponent(component)
					txctx.Logger.WithError(e).Errorf("nonce: could not load recovery status")
					return
				}
				if !recovering {
					txctx.Logger.WithFields(log.Fields{
						"nonce.recovering.expected": expectedNonce,
					}).Warnf("nonce: start recovering")
					// Envelope being processed is the first one we encounter with Nonce too High
					// Indicate expected nonce in metadata to signal tx-nonce worker from which value re-start attributing nonces
					txctx.Envelope.SetMetadataValue("nonce.recovering.expected", strconv.FormatUint(expectedNonce, 10))
				}
			}

			// We set a context value to indicate to other handlers that
			// an invalid nonce has been processed
			txctx.Set("invalid.nonce", true)

			// Abort execution
			txctx.Abort()

			return
		}

		// Execute pending handlers
		txctx.Next()

		// If pending handlers executed correctly we increment nonce
		if len(txctx.Envelope.GetErrors()) == 0 {
			// Increment last sent nonce
			err := nm.SetLastSent(nonceKey, n)
			if err != nil {
				// An error here means that we probably lost connection with NonceManager underlying cache.
				// TODO: A retry strategy should be implemented on the nonce manager to make this scenario rare
				//
				// At this point pending handlers have been executed correctly
				// meaning that transaction has been sent to the chain, so we log
				// error and proceed to next envelope
				e := errors.FromError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Errorf("nonce: could not set nonce on cache")
			}
		}
	}
}

// RecoveryStatusSetter returns and handler responsible to set nonce recovery status after envelope has been processed
//
// We set recovery status in a separated middleware handler that should be surrounding producer handler
// So in case of nonce too high, we make sure that envelope has been effectively produced in tx-nonce topic before updating
// recovery status on nonce manager.
//
// Setting recovery status before producing the envelope could result in case of crash
// in a situation were would never be able to signal tx-nonce to recalibrate nonce
func RecoveryStatusSetter(nm nonce.Sender) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// Execute pending handlers
		txctx.Next()

		// Retrieves nonce key for nonce manager procesing
		nonceKey := string(txctx.In.Key())

		// Test if we have entered recovery mode
		// In which case we set recovery status to true
		if _, ok := txctx.Envelope.GetMetadataValue("nonce.recovering.expected"); ok {
			// We encountered a first nonce too high scenario so we set recovery status to true
			err := nm.SetRecovering(nonceKey, true)
			if err != nil {
				txctx.Logger.WithError(
					errors.FromError(err).ExtendComponent(component),
				).Errorf("nonce: could not load recovery status")
				return
			}
		}

		if b, ok := txctx.Get("invalid.nonce").(bool); len(txctx.Envelope.GetErrors()) == 0 && (!ok || !b) {
			// Transaction has been processed properly
			// Deactivate recovery if activated
			recovering, err := nm.IsRecovering(nonceKey)
			if err != nil {
				txctx.Logger.WithError(
					errors.FromError(err).ExtendComponent(component),
				).Errorf("nonce: could not load recovery status from cache")
			}
			if recovering {
				// Envelope being processed is valid meaning nonce recovery
				// has completed. So we deactivate recovery status
				err := nm.SetRecovering(nonceKey, false)
				if err != nil {
					txctx.Logger.WithError(
						errors.FromError(err).ExtendComponent(component),
					).Errorf("nonce: could not set recovery status from cache")
				}
			}
		}
	}
}

func resetTx(tx *ethereum.Transaction) {
	if tx.GetTxData() != nil {
		tx.GetTxData().SetNonce(0)
	}
	tx.Hash = nil
	tx.Raw = nil
}
