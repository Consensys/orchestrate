package sender

import (
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	contextStore "gitlab.com/ConsenSys/client/fr/core-stack/service/envelope-store.git/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/types"
)

// Sender creates a Sender handler
func Sender(sender ethclient.TransactionSender, store contextStore.EnvelopeStore) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"chain.id": txctx.Envelope.GetChain().GetId(),
		})

		if isUnsignedTx(txctx) {
			proccessTxWithUndeterministicHash(txctx, store, func() (common.Hash, error) {
				return sendPublicUnsignedTx(txctx, sender)
			})
			return
		}

		processTxWithDeterministicHash(txctx, store, func() error {
			return sender.SendRawTransaction(txctx.Context(), txctx.Envelope.GetChain().ID(), txctx.Envelope.GetTx().GetRaw().Hex())
		})
	}
}

func isUnsignedTx(txctx *engine.TxContext) bool {
	return txctx.Envelope.GetTx().GetRaw().GetRaw() == nil || txctx.Envelope.GetTx().GetHash().GetRaw() == nil
}

func proccessTxWithUndeterministicHash(txctx *engine.TxContext, store contextStore.EnvelopeStore, executeTx func() (common.Hash, error)) {
	txHash, err := executeTx()
	if err != nil {
		txctx.Logger.WithError(err).Errorf("sender: could not send transaction")
		_ = txctx.AbortWithError(err)
		return
	}

	// Set transaction Hash on trace
	txctx.Envelope.GetTx().SetHash(txHash)
	txctx.Logger = txctx.Logger.WithFields(log.Fields{
		"tx.hash": txctx.Envelope.GetTx().GetHash(),
	})
	txctx.Logger.Debugf("sender: transaction sent")

	// Store trace
	// We can not store trace before sending transaction because we do not know the transaction hash
	// This is an issue for overall consistency of the system before/after transaction is mined
	txctx.Logger.Infof("%v %v %v", txctx.Envelope.Chain.Id, txctx.Envelope.Tx.Hash, txctx.Envelope.Metadata.Id)
	_, _, err = store.Store(txctx.Context(), txctx.Envelope)
	if err != nil {
		// Connection to store is broken
		txctx.Logger.WithError(err).Errorf("sender: trace store failed to store trace")
		_ = txctx.AbortWithError(err)
		return
	}

	// Transaction has been properly sent so we set status to `pending`
	err = store.SetStatus(txctx.Context(), txctx.Envelope.GetMetadata().GetId(), "pending")
	if err != nil {
		// Connection to store is broken
		txctx.Logger.WithError(err).Errorf("sender: piou trace store failed to set status")
		_ = txctx.Error(err)
		return
	}
}

func sendPublicUnsignedTx(txctx *engine.TxContext, sender ethclient.TransactionSender) (common.Hash, error) {
	args := types.Envelope2SendTxArgs(txctx.Envelope)
	txHash, err := sender.SendTransaction(txctx.Context(), txctx.Envelope.GetChain().ID(), args)
	return txHash, err
}

func processTxWithDeterministicHash(txctx *engine.TxContext, store contextStore.EnvelopeStore, executeTx func() error) {
	txctx.Logger = txctx.Logger.WithFields(log.Fields{
		"tx.raw":  txctx.Envelope.GetTx().GetRaw().Hex(),
		"tx.hash": txctx.Envelope.GetTx().GetHash(),
	})

	log.WithFields(log.Fields{
		"nonce": txctx.Envelope.GetTx().GetTxData().GetNonce(),
		"from":  txctx.Envelope.GetFrom().Hex(),
	}).Info("processing transaction")

	// Store trace
	status, _, err := store.Store(txctx.Context(), txctx.Envelope)
	if err != nil {
		// Connection to store is broken
		txctx.Logger.WithError(err).Errorf("sender: trace store failed to store trace")
		_ = txctx.AbortWithError(err)
		return
	}

	if status == "pending" {
		// Tx has already been sent
		// TODO: Request TxHash from chain to make sure we do not miss a message
		txctx.Logger.Warnf("sender: transaction has already been sent")
		txctx.Abort()
		return
	}

	err = executeTx()
	if err != nil {
		txctx.Logger.WithError(err).Errorf("sender: could not send transaction")

		// TODO: handle error
		_ = txctx.Error(err)

		// We update status in storage
		storeErr := store.SetStatus(txctx.Context(), txctx.Envelope.GetMetadata().GetId(), "error")
		if storeErr != nil {
			// Connection to store is broken
			txctx.Logger.WithError(storeErr).Errorf("sender: trace store failed to set status")
			_ = txctx.Error(storeErr)
		}
		txctx.Abort()
		return
	}
	txctx.Logger.Debugf("sender: raw transaction sent")

	// Transaction has been properly sent so we set status to `pending`
	err = store.SetStatus(txctx.Context(), txctx.Envelope.GetMetadata().GetId(), "pending")
	if err != nil {
		// Connection to store is broken
		txctx.Logger.WithError(err).Errorf("sender: trace store failed to set status")
		_ = txctx.Error(err)
		return
	}
}
