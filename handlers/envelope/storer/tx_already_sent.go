package storer

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/proxy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/client"
)

// TxAlreadySent implements an handler that controls whether transaction associated to envelope
// has already been sent and abort execution of pending handlers
//
// This handler makes guarantee that envelopes with the same UUID will not be send twice (scenario that could append in case
// of crash. As transaction orchestration system is configured to consume Kafka messages at least once).
func TxAlreadySent(ec ethclient.ChainLedgerReader, txSchedulerClient client.TransactionSchedulerClient) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Logger.Tracef("from TxAlreadySent => TenantID value: %s", multitenancy.TenantIDFromContext(txctx.Context()))

		// Load possibly already sent envelope
		job, err := txSchedulerClient.GetJob(txctx.Context(), txctx.Envelope.GetJobUUID())
		if err != nil && !errors.IsNotFoundError(err) {
			// Connection to tx scheduler is broken
			e := txctx.AbortWithError(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Error("transaction scheduler: failed to get job")
			return
		}

		// Tx has already been updated
		switch job.Status {
		case utils.StatusPending:
			if txctx.Envelope.IsResendingJobTx() {
				txctx.Logger.Debug("transaction scheduler: transaction is being resent")
				return
			}

			txctx.SetTxAlreadyPending(true)
			txctx.Logger.Warn("transaction scheduler: transaction has already been updated")

			url, err := proxy.GetURL(txctx)
			if err != nil {
				return
			}

			// We make sure that transaction has already been sent to the ETH node by querying to chain
			tx, _, err := ec.TransactionByHash(txctx.Context(), url, job.Transaction.GetHash())
			if err != nil {
				e := txctx.AbortWithError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Error("transaction scheduler: connection to Ethereum client is broken")
			}

			if tx != nil {
				// Transaction has already been sent so we abort execution
				txctx.Logger.Warn("transaction scheduler: transaction has already been sent but status was not set")
				txctx.Abort()
			}
		case utils.StatusMined:
			// Transaction has already been sent so we abort execution
			txctx.Logger.Warn("transaction scheduler: transaction has already been sent")
			txctx.Abort()
		case utils.StatusFailed:
			// Transaction has already been failed so we abort execution
			txctx.Logger.Warn("transaction scheduler: transaction has already been failed")
			txctx.Abort()
		default:
			var txHash string
			if txctx.Envelope.TxHash != nil {
				txHash = txctx.Envelope.TxHash.String()
			}

			txctx.Logger.WithField("txHash", txHash).Debug("transaction scheduler: transaction has not been sent")
		}
	}
}
