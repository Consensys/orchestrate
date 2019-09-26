package store

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
	evlpstore "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services/envelope-store"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/ethclient"
)

// TxAlreadySent implements an handler that controls wether transaction associated to envelope
// has already been sent and abort execution of pending handlers
//
// This handler makes guarantee that envelopes with the same ID will not be send twice (scenario that could append in case
// of crash. As transaction orchestration system is configured to consume Kafka messages at least once).
//
// Warning: above guarantee require embedded handler to
// 1. Store envelope on Envelope store
// 2. Send transaction to blockchain
// 3. Set envelope status
func TxAlreadySent(ec ethclient.ChainLedgerReader, s evlpstore.EnvelopeStoreClient) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// Load possibly already sent envelope
		resp, err := s.LoadByID(
			txctx.Context(),
			&evlpstore.LoadByIDRequest{
				Id: txctx.Envelope.GetMetadata().GetId(),
			})
		if err != nil && !errors.IsNotFoundError(err) {
			// Connection to store is broken
			e := txctx.AbortWithError(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("store: envelope store failed to store envelope")
			return
		}

		// Tx has already been stored
		if resp.GetStatusInfo().HasBeenSent() {
			txctx.Logger.Warnf("store: transaction has already been stored")

			// We make sure that transaction has not already been sent
			// by querying the chain
			tx, _, err := ec.TransactionByHash(
				txctx.Context(),
				resp.GetEnvelope().GetChain().ID(),
				resp.GetEnvelope().GetTx().GetHash().Hash(),
			)

			if err != nil {
				// Connection to Ethereum node is broken
				e := txctx.AbortWithError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Errorf("store: connection to Ethereum client is broken")
				return
			}

			if tx != nil {
				// Transaction has already been sent so we abort execution
				txctx.Logger.Warnf("store: transaction has already been sent")
				txctx.Abort()
				return
			}
		}

		txctx.Logger.Debugf("store: transaction has not been sent")
	}
}
