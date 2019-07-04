package sender

import (
	"errors"
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/utils"

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

		if isPublicTx(txctx) {
			processPublicTx(txctx, sender, store)
			return
		}

		processPrivateTx(txctx, sender, store)
	}
}

func isPublicTx(txctx *engine.TxContext) bool {
	return txctx.Envelope.GetArgs().GetPrivate() == nil
}

func processPublicTx(txctx *engine.TxContext, sender ethclient.TransactionSender, store contextStore.EnvelopeStore) {
	if !txctx.Envelope.GetTx().IsSigned() {
		processTxWithNonDeterministicHash(txctx, store, func() (common.Hash, error) {
			return sendPublicUnsignedTx(txctx, sender)
		})
		return
	}

	processTxWithDeterministicHash(txctx, store, func() error {
		return sender.SendRawTransaction(txctx.Context(), txctx.Envelope.GetChain().ID(), txctx.Envelope.GetTx().GetRaw().Hex())
	})
}

func processPrivateTx(txctx *engine.TxContext, sender ethclient.TransactionSender, store contextStore.EnvelopeStore) {

	protocol := txctx.Envelope.GetProtocol()
	switch {
	case protocol == nil:
		err := errors.New("protocol should be specified to send a private transaction")
		txctx.Logger.WithError(err).Errorf("sender: could not send private transaction")
		_ = txctx.AbortWithError(err)

	case protocol.IsPantheon():
		processTxWithNonDeterministicHash(txctx, store, func() (common.Hash, error) {
			return sendRawPrivateTx(txctx, sender)
		})

	case protocol.IsTessera():
		processTxWithDeterministicHash(txctx, store, func() error {
			return sendRawQuorumPrivateTx(txctx, sender)
		})

	// It is not Tessera, but it can be Quorum with Constellation
	case protocol.IsConstellation():
		if txctx.Envelope.GetTx().IsSigned() {
			err := errors.New("transactions executed with Constellation should be unsigned")
			txctx.Logger.WithError(err).Errorf("sender: could not send private transaction")
			_ = txctx.AbortWithError(err)
		}
		processTxWithNonDeterministicHash(txctx, store, func() (common.Hash, error) {
			return sendPublicUnsignedTx(txctx, sender)
		})
	default:
		err := fmt.Errorf("cannot process a private transaction with protocol %s", protocol.String())
		txctx.Logger.WithError(err).Errorf("sender: could not send private transaction")
		_ = txctx.AbortWithError(err)
	}
}

func sendRawPrivateTx(txctx *engine.TxContext, sender ethclient.TransactionSender) (common.Hash, error) {
	txctx.Logger.Infof("sender: sending raw private transaction")
	privateArgs := types.Call2PrivateArgs(txctx.Envelope.GetArgs())
	return sender.SendRawPrivateTransaction(txctx.Context(), txctx.Envelope.GetChain().ID(), txctx.Envelope.GetTx().GetRaw().GetRaw(), privateArgs)
}

func sendRawQuorumPrivateTx(txctx *engine.TxContext, sender ethclient.TransactionSender) error {
	txctx.Logger.Infof("sender: sending raw Quorum private transaction")
	privateArgs := types.Call2PrivateArgs(txctx.Envelope.GetArgs())
	hash, err := sender.SendQuorumRawPrivateTransaction(txctx.Context(), txctx.Envelope.GetChain().ID(), txctx.Envelope.GetTx().GetRaw().GetRaw(), privateArgs.PrivateFor)

	if err == nil {
		txctx.Logger.Infof("sender: result transaction hash is %s", hash.Hex())
	}

	return err
}

func processTxWithNonDeterministicHash(txctx *engine.TxContext, store contextStore.EnvelopeStore, sendTx func() (common.Hash, error)) {
	txHash, err := sendTx()
	if err != nil {
		txctx.Logger.WithError(err).Errorf("sender: could not send transaction")
		_ = txctx.AbortWithError(err)
		return
	}

	// Set transaction Hash on envelope
	txctx.Envelope.GetTx().SetHash(txHash)
	txctx.Logger = txctx.Logger.WithFields(log.Fields{
		"tx.hash": txctx.Envelope.GetTx().GetHash(),
	})
	txctx.Logger.Debugf("sender: transaction sent")

	// Store envelope
	// We can not store envelope before sending transaction because we do not know the transaction hash
	// This is an issue for overall consistency of the system before/after transaction is mined
	txctx.Logger.Infof("%v %v %v", txctx.Envelope.Chain.Id, txctx.Envelope.Tx.Hash, txctx.Envelope.Metadata.Id)
	_, _, err = store.Store(txctx.Context(), txctx.Envelope)
	if err != nil {
		// Connection to store is broken
		txctx.Logger.WithError(err).Errorf("sender: envelope store failed to store envelope")
		_ = txctx.AbortWithError(err)
		return
	}

	// Transaction has been properly sent so we set status to `pending`
	err = store.SetStatus(txctx.Context(), txctx.Envelope.GetMetadata().GetId(), "pending")
	if err != nil {
		// Connection to store is broken
		txctx.Logger.WithError(err).Errorf("sender: envelope store failed to set status")
		_ = txctx.Error(err)
		return
	}
}

func sendPublicUnsignedTx(txctx *engine.TxContext, sender ethclient.TransactionSender) (common.Hash, error) {
	txctx.Logger.Infof("sender: sending public unsigned transaction")
	args := types.Envelope2SendTxArgs(txctx.Envelope)
	txHash, err := sender.SendTransaction(txctx.Context(), txctx.Envelope.GetChain().ID(), args)
	return txHash, err
}

func processTxWithDeterministicHash(txctx *engine.TxContext, store contextStore.EnvelopeStore, sendTx func() error) {
	txctx.Logger = txctx.Logger.WithFields(log.Fields{
		"tx.raw":  utils.ShortString(txctx.Envelope.GetTx().GetRaw().Hex(), 30),
		"tx.hash": txctx.Envelope.GetTx().GetHash(),
	})

	log.WithFields(log.Fields{
		"nonce": txctx.Envelope.GetTx().GetTxData().GetNonce(),
		"from":  txctx.Envelope.GetFrom().Hex(),
	}).Info("processing transaction")

	// Store envelope
	status, _, err := store.Store(txctx.Context(), txctx.Envelope)
	if err != nil {
		// Connection to store is broken
		txctx.Logger.WithError(err).Errorf("sender: envelope store failed to store envelope")
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

	err = sendTx()
	if err != nil {
		txctx.Logger.WithError(err).Errorf("sender: could not send transaction")

		// TODO: handle error
		_ = txctx.Error(err)

		// We update status in storage
		storeErr := store.SetStatus(txctx.Context(), txctx.Envelope.GetMetadata().GetId(), "error")
		if storeErr != nil {
			// Connection to store is broken
			txctx.Logger.WithError(storeErr).Errorf("sender: envelope store failed to set envelope")
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
		txctx.Logger.WithError(err).Errorf("sender: envelope store failed to set status")
		_ = txctx.Error(err)
		return
	}
}
