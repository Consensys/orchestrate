package noncechecker

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/nonce/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient"
	ierror "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/error"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/nonce"
)

func controlRecoveryCount(txctx *engine.TxContext, conf *Configuration) error {
	var count int
	if v := txctx.Envelope.GetInternalLabelsValue("nonce.recovering.count"); v != "" {
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
	_ = txctx.Envelope.SetInternalLabelsValue("nonce.recovering.count", fmt.Sprintf("%v", count+1))
	return nil
}

// Checker creates an handler responsible to check transaction nonce value
// It makes sure that transactions with invalid nonce are not processed
func Checker(conf *Configuration, nm nonce.Sender, ec ethclient.ChainStateReader, tracker *RecoveryTracker) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		if mode := txctx.Envelope.GetContextLabelsValue("txMode"); mode == "raw" {
			// If transaction has been generated externally we skip nonce check
			txctx.Logger.WithFields(log.Fields{
				"id": txctx.Envelope.GetID(),
			}).Debugf("nonce: skip check for raw transaction")
			return
		}

		if txctx.Envelope.IsOneTimeKeySignature() {
			// If transaction has been generated externally we skip nonce check
			txctx.Logger.WithFields(log.Fields{
				"id": txctx.Envelope.GetID(),
			}).Debugf("nonce: skip check for one-time-key signing")
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
			"id":      txctx.Envelope.GetID(),
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

		url, err := proxy.GetURL(txctx)
		if err != nil {
			return
		}

		if ok {
			expectedNonce = lastSent + 1
		} else {
			// If no nonce is available (meaning that envelope being processed is the first one for the pair sender,chain)
			// then we retrieve nonce from chain
			pendingNonce, callErr := utils.GetNonce(ec, txctx, url)
			if callErr != nil {
				_ = txctx.AbortWithError(errors.EthereumError("could not read nonce from chain - got %v", callErr)).ExtendComponent(component)
				return
			}
			txctx.Logger.WithFields(log.Fields{
				"nonce.pending": pendingNonce,
			}).Debugf("nonce: calibrating nonce from chain")

			expectedNonce = pendingNonce
		}

		n, err := txctx.Envelope.GetNonceUint64()
		if err != nil {
			_ = txctx.AbortWithError(errors.DataError("could not check nonce - %s", err)).ExtendComponent(component)
			return
		}
		if n != expectedNonce {
			// Attributes nonce is invalid
			txctx.Logger.WithFields(log.Fields{
				"nonce.expected": expectedNonce,
				"nonce.got":      n,
			}).Warnf("nonce: invalid nonce")

			// Reset tx nonce, hash and raw
			resetTx(txctx.Envelope)

			// Control recovery count
			err := controlRecoveryCount(txctx, conf)
			if err != nil {
				e := txctx.AbortWithError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Errorf("nonce: abort transaction execution")
				return
			}

			if n > expectedNonce {
				// Nonce is too-high
				if nonceRecoveringExpected := txctx.Envelope.GetInternalLabelsValue("nonce.recovering.expected"); nonceRecoveringExpected != "" || tracker.Recovering(nonceKey) == 0 {
					// Envelope has already been recovered or it is the first one to be recovered
					_ = txctx.Envelope.SetInternalLabelsValue("nonce.recovering.expected", strconv.FormatUint(expectedNonce, 10))
				}
			} else {
				// If nonce is to low we remove any recovery signal in metadata (possibly coming from a prior execution)
				delete(txctx.Envelope.InternalLabels, "nonce.recovering.expected")
			}

			// We set a context value to indicate to other handlers that
			// an invalid nonce has been processed
			txctx.Set("invalid.nonce", true)

			// Abort execution
			txctx.Abort()

			return
		}

		// If nonce was valid we remove recovery signal in metadata (possibly coming from a prior execution)
		delete(txctx.Envelope.InternalLabels, "nonce.recovering.expected")
		delete(txctx.Envelope.InternalLabels, "nonce.recovering.count")

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
			if !strings.Contains(err.String(), "nonce too low") && !strings.Contains(err.String(), "Nonce too low") && !strings.Contains(err.String(), "Incorrect nonce") {
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
				_ = txctx.AbortWithError(errors.FromError(fmt.Errorf("abort transaction execution - got %v", err))).ExtendComponent(component)
				return
			}

			// We recalibrate nonce from chain
			pendingNonce, err := utils.GetNonce(ec, txctx, url)
			if err != nil {
				_ = txctx.AbortWithError(errors.FromError(fmt.Errorf("could not read nonce from chain - got %v", err))).ExtendComponent(component)
				return
			}
			txctx.Logger.WithFields(log.Fields{
				"nonce.pending": pendingNonce,
			}).Debugf("nonce: calibrating nonce from chain")

			// Re-calibrate cache
			err = nm.SetLastSent(nonceKey, pendingNonce-1)
			if err != nil {
				_ = txctx.AbortWithError(errors.FromError(fmt.Errorf("could not store last sent nonce - got %v", err))).ExtendComponent(component)
			}

			// Reset tx nonce, hash and raw
			resetTx(txctx.Envelope)

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
// So in case of nonce too high, we make sure that envelope has been effectively produced in tx-crafter topic before updating
// recovery status on nonce manager.
//
// Setting recovery status before producing the envelope could result in case of crash
// in a situation were would never be able to signal tx-crafter to recalibrate nonce
func RecoveryStatusSetter(nm nonce.Sender, tracker *RecoveryTracker) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// Execute pending handlers
		txctx.Next()

		// Compute nonce key
		nonceKey := string(txctx.In.Key())

		// If are recovering we increment count
		if nonceRecoveringExpected := txctx.Envelope.GetInternalLabelsValue("nonce.recovering.expected"); nonceRecoveringExpected != "" {
			tracker.Recover(nonceKey)
		}

		if b, ok := txctx.Get("invalid.nonce").(bool); !ok || !b {
			// Transaction has been processed properly
			// Deactivate recovery if activated
			tracker.Recovered(nonceKey)
		}
	}
}

func resetTx(req *tx.Envelope) {
	req.Nonce = nil
	req.TxHash = nil
	req.Raw = ""
}
