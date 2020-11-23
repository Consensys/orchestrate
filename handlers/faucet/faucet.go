package faucet

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	multitenancy2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/chainregistry"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/txscheduler"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store/models"
	client2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/client"
)

const component = "handler.faucet"

// Faucet creates a Faucet handler
func Faucet(registryClient client.ChainRegistryClient, txSchedulerClient client2.TransactionSchedulerClient) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		logger := txctx.Logger.
			WithField("envelope_id", txctx.Envelope.GetID()).
			WithField("job_uuid", txctx.Envelope.GetJobUUID()).
			WithField("tenant_id", multitenancy2.TenantIDFromContext(txctx.Context())).
			WithField("account", txctx.Envelope.MustGetFromAddress().Hex())

		logger.Debugf("faucet handler starts")

		var chain *models.Chain
		var err error
		switch {
		case txctx.Envelope.GetChainName() != "":
			chain, err = registryClient.GetChainByName(txctx.Context(), txctx.Envelope.GetChainName())
			if err != nil {
				_ = txctx.Error(errors.FaucetWarning("could not find chain name %s: %v", txctx.Envelope.GetChainName(), err)).ExtendComponent(component)
				return
			}
		case txctx.Envelope.GetChainUUID() != "":
			chain, err = registryClient.GetChainByUUID(txctx.Context(), txctx.Envelope.GetChainUUID())
			if err != nil {
				_ = txctx.Error(errors.FaucetWarning("could not find chain uuid %s: %v", txctx.Envelope.GetChainUUID(), err)).ExtendComponent(component)
				return
			}
		default:
			err := errors.FaucetWarning("skipped because no chain attached to envelope").ExtendComponent(component)
			logger.Debugf(err.Error())
			return
		}

		logger = txctx.Logger.WithField("chain_uuid", chain.UUID)

		fct, err := registryClient.GetFaucetCandidate(txctx.Context(), txctx.Envelope.MustGetFromAddress(), chain.UUID)
		if err != nil {
			if errors.IsNotFoundError(err) {
				logger.WithError(err).Debugf("could not get any faucet candidate")
				return
			}

			logger.WithError(err).Error("failed to fetch faucet candidate")
			_ = txctx.Error(err).ExtendComponent(component)
			return
		}

		_, err = txSchedulerClient.SendTransferTransaction(txctx.Context(),
			&types.TransferRequest{
				ChainName: chain.Name,
				Params: types.TransferParams{
					From:  fct.Creditor.Hex(),
					To:    txctx.Envelope.MustGetFromAddress().String(),
					Value: fct.Amount.String(),
				},
				Labels: chainregistry.FaucetToJobLabels(fct),
			})
		if err != nil {
			logger.WithError(err).Error("fail to send funding transaction")
			_ = txctx.Error(err).ExtendComponent(component)
			return
		}

		logger.WithFields(log.Fields{
			"faucet.amount": fct.Amount.Text(10),
		}).Infof("faucet: credit approved")
	}
}
