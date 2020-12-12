package units

import (
	"context"
	"math/rand"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/tx"
	txschedulertypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/txscheduler"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	txscheduler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/client"
	utils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/service/stress/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/utils/chanregistry"
)

func BatchDeployContractTest(ctx context.Context, cfg *WorkloadConfig, client txscheduler.TransactionSchedulerClient, chanReg *chanregistry.ChanRegistry) error {
	accounts := utils2.ContextAccounts(ctx)
	chains := utils2.ContextChains(ctx)
	log.FromContext(ctx).Debugf("Running batchDeployContract()...")

	nAcc := rand.Intn(cfg.nAccounts)
	idempotency := utils.RandomString(30)
	evlp := tx.NewEnvelope()
	t := utils2.NewEnvelopeTracker(chanReg, evlp, idempotency)

	_, err := client.SendDeployTransaction(ctx, &txschedulertypes.DeployContractRequest{
		ChainName: chains["besu"].Name,
		Params: txschedulertypes.DeployContractParams{
			From:         accounts[nAcc],
			ContractName: "SimpleToken",
		},
		Labels: map[string]string{
			"id": idempotency,
		},
	})

	if err != nil {
		log.FromContext(ctx).WithError(err).Error("failed to send transaction")
		return err
	}

	err = t.Load("tx.decoded", cfg.waitForEnvelopeTimeout)
	if err != nil {
		log.FromContext(ctx).WithError(err).Error("failed to fetch envelope")
		return err
	}

	log.FromContext(ctx).Debugf("Done: Envelope was found in tx-decoded")
	return nil
}
