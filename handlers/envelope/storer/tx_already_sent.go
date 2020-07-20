package storer

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
)

// TxAlreadySent implements an handler that controls whether transaction associated to envelope
// has already been sent and abort execution of pending handlers
//
// This handler makes guarantee that envelopes with the same UUID will not be send twice (scenario that could append in case
// of crash. As transaction orchestration system is configured to consume Kafka messages at least once).
//
// Warning: above guarantee require embedded handler to
// 1. Store envelope on Envelope store
// 2. Send transaction to blockchain
// 3. Set envelope status
func TxAlreadySent(ec ethclient.ChainLedgerReader, s svc.EnvelopeStoreClient, txSchedulerClient client.TransactionSchedulerClient) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// TODO: Remove this if when envelope store is removed
		if txctx.Envelope.BelongToEnvelopeStore() {
			checkTxInStore(txctx, ec, s)
		} else {
			checkTxInScheduler(txctx, ec, txSchedulerClient)
		}
	}
}

func checkTxInScheduler(txctx *engine.TxContext, ec ethclient.ChainLedgerReader, txSchedulerClient client.TransactionSchedulerClient) {
	txctx.Logger.Tracef("from TxAlreadySent => TenantID value: %s", multitenancy.TenantIDFromContext(txctx.Context()))

	// Load possibly already sent envelope
	job, err := txSchedulerClient.GetJob(txctx.Context(), txctx.Envelope.GetID())
	if err != nil && !errors.IsNotFoundError(err) {
		// Connection to tx scheduler is broken
		e := txctx.AbortWithError(err).ExtendComponent(component)
		txctx.Logger.WithError(e).Errorf("transaction scheduler: failed to get job")
		return
	}

	// Tx has already been updated
	if job.Status == utils.StatusPending {
		txctx.Logger.Warnf("transaction scheduler: transaction has already been updated")
		url, err := proxy.GetURL(txctx)
		if err != nil {
			return
		}

		// We make sure that transaction has not already been sent to the ETH node by querying to chain
		tx, _, err := ec.TransactionByHash(txctx.Context(), url, job.Transaction.GetHash())
		if err != nil {
			// Connection to Ethereum node is broken
			e := txctx.AbortWithError(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("transaction scheduler: connection to Ethereum client is broken")
			return
		}

		if tx != nil {
			// Transaction has already been sent so we abort execution
			txctx.Logger.Warnf("transaction scheduler: transaction has already been sent but status was not set")
			txctx.Abort()
			return
		}
	} else if job.Status == utils.StatusRecovering || job.Status == utils.StatusMined {
		// Transaction has already been sent so we abort execution
		txctx.Logger.Warnf("transaction scheduler: transaction has already been sent")
		txctx.Abort()
		return
	}

	txctx.Logger.
		WithField("txHash", txctx.Envelope.TxHash.String()).
		Debugf("transaction scheduler: transaction has not been sent")
}

func checkTxInStore(txctx *engine.TxContext, ec ethclient.ChainLedgerReader, s svc.EnvelopeStoreClient) {
	tenantID := multitenancy.TenantIDFromContext(txctx.Context())
	txctx.Logger.Tracef("from TxAlreadySent => TenantID value: %s", tenantID)
	// Load possibly already sent envelope
	resp, err := s.LoadByID(
		txctx.Context(),
		&svc.LoadByIDRequest{
			Id: txctx.Envelope.GetID(),
		},
	)
	if err != nil && !errors.IsNotFoundError(err) {
		// Connection to store is broken
		e := txctx.AbortWithError(err).ExtendComponent(component)
		txctx.Logger.WithError(e).Errorf("store: envelope store failed to store envelope")
		return
	}

	// Tx has already been stored
	if resp.GetStatusInfo().HasBeenSent() {
		txctx.Logger.Warnf("store: transaction has already been stored")
		url, err := proxy.GetURL(txctx)
		if err != nil {
			return
		}

		// We make sure that transaction has not already been sent
		// by querying the chain
		tx, _, err := ec.TransactionByHash(
			txctx.Context(),
			url,
			resp.GetEnvelope().TxHash(),
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
