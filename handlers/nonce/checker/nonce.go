package noncechecker

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/nonce"
	ierror "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/error"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/ethereum"
)

func controlRecoveryCount(txctx *engine.TxContext, conf *Configuration) error {
	var count int
	v, ok := txctx.Envelope.GetMetadataValue("nonce.recovering.count")
	if ok {
		i, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		count = i
	}

	if count >= conf.MaxRecovery {
		// If maximum recovery count is reached do not recover
		return fmt.Errorf("nonce: reached max recovery count")
	}

	// Incremenent recovery count on envelope
	txctx.Envelope.SetMetadataValue("nonce.recovering.count", fmt.Sprintf("%v", count+1))
	return nil
}

// Checker creates an handler responsible to check transaction nonce value
// It makes sure that transactions with invalid nonce are not processed
func Checker(conf *Configuration, nm nonce.Sender, ec ethclient.ChainStateReader, tracker *RecoveryTracker) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		if mode, ok := txctx.Envelope.GetMetadataValue("tx.mode"); ok && mode == "raw" {
			// If transaction has been generated externally we skip nonce check
			txctx.Logger.WithFields(log.Fields{
				"metadata.id": txctx.Envelope.GetMetadata().GetId(),
			}).Debugf("nonce: skip check for raw transaction")
			return
		}

		// Retrieve chainID and sender address
		chainID, sender := txctx.Envelope.GetChain().ID(), txctx.Envelope.GetFrom().Address()
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"tx.sender":   sender.Hex(),
			"chain.id":    chainID.String(),
			"metadata.id": txctx.Envelope.GetMetadata().GetId(),
		})

		// Retrieves nonce key for nonce manager processing
		nonceKey := string(txctx.In.Key())
		var expectedNonce uint64

		// Retrieve last sent nonce from nonce manager
		lastSent, ok, err := nm.GetLastSent(nonceKey)
		if err != nil {
			e := txctx.AbortWithError(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("nonce: could not load last sent nonce")
			return
		}

		if ok {
			expectedNonce = lastSent + 1
		} else {
			// If no nonce is available (meaning that envelope being processed is the first one for the pair sender,chain)
			// then we retrieve nonce from chain
			pendingNonce, err := ec.PendingNonceAt(txctx.Context(), chainID, sender)
			if err != nil {
				e := txctx.AbortWithError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Errorf("nonce: could not read nonce from chain")
				return
			}
			txctx.Logger.WithFields(log.Fields{
				"nonce.pending": pendingNonce,
			}).Debugf("nonce: calibrating nonce from chain")

			expectedNonce = pendingNonce
		}

		n := txctx.Envelope.GetTx().GetTxData().GetNonce()
		if n != expectedNonce {
			// Attributes nonce is invalid
			txctx.Logger.WithFields(log.Fields{
				"nonce.expected": expectedNonce,
				"nonce.got":      n,
			}).Warnf("nonce: invalid nonce")

			// Reset tx nonce, hash and raw
			resetTx(txctx.Envelope.GetTx())

			// Control recovery count
			err := controlRecoveryCount(txctx, conf)
			if err != nil {
				e := txctx.AbortWithError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Errorf("nonce: abort transaction execution")
				return
			}

			if n > expectedNonce {
				// Nonce is too-high
				if _, ok := txctx.Envelope.GetMetadataValue("nonce.recovering.expected"); ok || tracker.Recovering(nonceKey) == 0 {
					// Envelope has already been recovered or it is the first one to be recovered
					txctx.Envelope.SetMetadataValue("nonce.recovering.expected", strconv.FormatUint(expectedNonce, 10))
				}
			} else {
				// If nonce is to low we remove any recovery signal in metadata (possibly coming from a prior execution)
				delete(txctx.Envelope.GetMetadata().GetExtra(), "nonce.recovering.expected")
			}

			// We set a context value to indicate to other handlers that
			// an invalid nonce has been processed
			txctx.Set("invalid.nonce", true)

			// Abort execution
			txctx.Abort()

			return
		}

		// If nonce was valid we remove recovery signal in metadata (possibly coming from a prior execution)
		delete(txctx.Envelope.GetMetadata().GetExtra(), "nonce.recovering.expected")
		delete(txctx.Envelope.GetMetadata().GetExtra(), "nonce.recovering.count")

		// Execute pending handlers
		txctx.Next()

		// If pending handlers executed correctly we increment nonce
		if len(txctx.Envelope.GetErrors()) == 0 {
			// Possibly re-initiliaze recovery counter
			tracker.Recovered(nonceKey)

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
				txctx.Logger.WithError(e).Errorf("nonce: could not store last sent nonce")
			}

			return
		}

		var errs []*ierror.Error
		for _, err := range txctx.Envelope.GetErrors() {
			// TODO: update EthClient to process and standardize nonce too low errors
			if !strings.Contains(err.String(), "nonce too low") {
				errs = append(errs, err)
				continue
			}

			// We got a "nonce too low" error when sending the transaction
			txctx.Logger.WithFields(log.Fields{
				"nonce.got": n,
			}).Warnf("nonce: invalid nonce")

			// Control recovery count
			err := controlRecoveryCount(txctx, conf)
			if err != nil {
				e := txctx.AbortWithError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Errorf("nonce: abort transaction execution")
				return
			}

			// We recalibrate nonce from chain
			pendingNonce, err := ec.PendingNonceAt(txctx.Context(), chainID, sender)
			if err != nil {
				e := txctx.AbortWithError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Errorf("nonce: could not read nonce from chain")
				return
			}
			txctx.Logger.WithFields(log.Fields{
				"nonce.pending": pendingNonce,
			}).Debugf("nonce: calibrating nonce from chain")

			// Re-calibrate cache
			err = nm.SetLastSent(nonceKey, pendingNonce-1)
			if err != nil {
				e := errors.FromError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Errorf("nonce: could not store last sent nonce")
			}

			// Reset tx nonce, hash and raw
			resetTx(txctx.Envelope.GetTx())

			// We set a context value to indicate to other handlers that
			// an invalid nonce has been processed
			txctx.Set("invalid.nonce", true)
		}

		// Update Envelope errors with Nonce too Low error removed
		txctx.Envelope.Errors = errs
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
func RecoveryStatusSetter(nm nonce.Sender, tracker *RecoveryTracker) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// Execute pending handlers
		txctx.Next()

		// Compute nonce key
		nonceKey := string(txctx.In.Key())

		// If are recovering we increment count
		if _, ok := txctx.Envelope.GetMetadataValue("nonce.recovering.expected"); ok {
			tracker.Recover(nonceKey)
		}

		if b, ok := txctx.Get("invalid.nonce").(bool); len(txctx.Envelope.GetErrors()) == 0 && (!ok || !b) {
			// Transaction has been processed properly
			// Deactivate recovery if activated
			tracker.Recovered(nonceKey)
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
