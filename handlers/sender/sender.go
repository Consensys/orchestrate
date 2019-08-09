package sender

import (
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
	evlpstore "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services/envelope-store"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/types"
)

// Sender creates a Sender handler
func Sender(sender ethclient.TransactionSender, store evlpstore.EnvelopeStoreClient) engine.HandlerFunc {
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

func processPublicTx(txctx *engine.TxContext, sender ethclient.TransactionSender, store evlpstore.EnvelopeStoreClient) {
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

func processPrivateTx(txctx *engine.TxContext, sender ethclient.TransactionSender, store evlpstore.EnvelopeStoreClient) {

	protocol := txctx.Envelope.GetProtocol()
	switch {
	case protocol == nil:
		err := errors.InvalidFormatError("protocol should be specified to send a private transaction").
			SetComponent(component)
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
			err := errors.InvalidFormatError("transactions executed with Constellation should be unsigned").
				SetComponent(component)
			txctx.Logger.WithError(err).Errorf("sender: could not send private transaction")
			_ = txctx.AbortWithError(err)
		}
		processTxWithNonDeterministicHash(txctx, store, func() (common.Hash, error) {
			return sendPublicUnsignedTx(txctx, sender)
		})
	default:
		err := errors.DataError("cannot process a private transaction with protocol %s", protocol.String()).
			SetComponent(component)
		txctx.Logger.WithError(err).Errorf("sender: could not send private transaction")
		_ = txctx.AbortWithError(err)
	}
}

func sendRawPrivateTx(txctx *engine.TxContext, sender ethclient.TransactionSender) (common.Hash, error) {
	txctx.Logger.Infof("sender: sending raw private transaction")
	privateArgs := types.Call2PrivateArgs(txctx.Envelope.GetArgs())
	hash, err := sender.SendRawPrivateTransaction(txctx.Context(), txctx.Envelope.GetChain().ID(), txctx.Envelope.GetTx().GetRaw().GetRaw(), privateArgs)
	if err != nil {
		return hash, errors.FromError(err).ExtendComponent(component)
	}
	return hash, nil
}

func sendRawQuorumPrivateTx(txctx *engine.TxContext, sender ethclient.TransactionSender) error {
	txctx.Logger.Infof("sender: sending raw Quorum private transaction")
	privateArgs := types.Call2PrivateArgs(txctx.Envelope.GetArgs())
	hash, err := sender.SendQuorumRawPrivateTransaction(txctx.Context(), txctx.Envelope.GetChain().ID(), txctx.Envelope.GetTx().GetRaw().GetRaw(), privateArgs.PrivateFor)
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	txctx.Logger.Infof("sender: result transaction hash is %s", hash.Hex())

	return err
}

func processTxWithNonDeterministicHash(txctx *engine.TxContext, store evlpstore.EnvelopeStoreClient, sendTx func() (common.Hash, error)) {
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
	_, err = store.Store(txctx.Context(), &evlpstore.StoreRequest{
		Envelope: txctx.Envelope,
	})
	if err != nil {
		// Connection to store is broken
		e := txctx.AbortWithError(err).ExtendComponent(component)
		txctx.Logger.WithError(e).Errorf("sender: envelope store failed to store envelope")
		return
	}

	// Transaction has been properly sent so we set status to `pending`
	_, err = store.SetStatus(txctx.Context(), &evlpstore.SetStatusRequest{
		Id:     txctx.Envelope.GetMetadata().GetId(),
		Status: evlpstore.Status_PENDING,
	})
	if err != nil {
		// Connection to store is broken
		e := errors.FromError(err).ExtendComponent(component)
		txctx.Logger.WithError(e).Warnf("sender: envelope store failed to set status")
		return
	}
}

func sendPublicUnsignedTx(txctx *engine.TxContext, sender ethclient.TransactionSender) (common.Hash, error) {
	txctx.Logger.Infof("sender: sending public unsigned transaction")
	args := types.Envelope2SendTxArgs(txctx.Envelope)
	txHash, err := sender.SendTransaction(txctx.Context(), txctx.Envelope.GetChain().ID(), args)
	if err != nil {
		return txHash, errors.FromError(err).ExtendComponent(component)
	}
	return txHash, nil
}

func processTxWithDeterministicHash(txctx *engine.TxContext, store evlpstore.EnvelopeStoreClient, sendTx func() error) {
	txctx.Logger = txctx.Logger.WithFields(log.Fields{
		"tx.raw":  utils.ShortString(txctx.Envelope.GetTx().GetRaw().Hex(), 30),
		"tx.hash": txctx.Envelope.GetTx().GetHash(),
	})

	log.WithFields(log.Fields{
		"nonce": txctx.Envelope.GetTx().GetTxData().GetNonce(),
		"from":  txctx.Envelope.GetFrom().Hex(),
	}).Info("processing transaction")

	// Store envelope
	resp, err := store.Store(
		txctx.Context(),
		&evlpstore.StoreRequest{
			Envelope: txctx.Envelope,
		})
	if err != nil {
		// Connection to store is broken
		e := txctx.AbortWithError(err).ExtendComponent(component)
		txctx.Logger.WithError(e).Errorf("sender: envelope store failed to store envelope")
		return
	}

	if resp.GetStatusInfo().GetStatus() == evlpstore.Status_PENDING {
		// Tx has already been sent
		// TODO: Request TxHash from chain to make sure we do not miss a message
		txctx.Logger.Warnf("sender: transaction has already been sent")
		txctx.Abort()
		return
	}

	err = sendTx()
	if err != nil {
		e := txctx.AbortWithError(err).ExtendComponent(component)
		txctx.Logger.WithError(e).Errorf("sender: could not send transaction")

		// We update status in storage
		_, storeErr := store.SetStatus(
			txctx.Context(),
			&evlpstore.SetStatusRequest{
				Id:     txctx.Envelope.GetMetadata().GetId(),
				Status: evlpstore.Status_ERROR,
			})
		if storeErr != nil {
			// Connection to store is broken
			e := errors.FromError(storeErr).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("sender: envelope store failed to set envelope")
		}
		return
	}
	txctx.Logger.Debugf("sender: raw transaction sent")

	// Transaction has been properly sent so we set status to `pending`
	_, err = store.SetStatus(
		txctx.Context(),
		&evlpstore.SetStatusRequest{
			Id:     txctx.Envelope.GetMetadata().GetId(),
			Status: evlpstore.Status_PENDING,
		})
	if err != nil {
		// Connection to store is broken
		e := errors.FromError(err).ExtendComponent(component)
		txctx.Logger.WithError(e).Errorf("sender: envelope store failed to set status")
		return
	}
}
