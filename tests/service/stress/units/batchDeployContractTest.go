package units

import (
	"context"
	"math/rand"
	"time"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	txschedulertypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx-scheduler"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	txscheduler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
	utils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/stress/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/utils/chanregistry"
)

var (
	timeout = time.Second * 10
)

func BatchDeployContractTest(ctx context.Context, client txscheduler.TransactionSchedulerClient, chanReg *chanregistry.ChanRegistry) error {
	accounts := utils2.ContextAccounts(ctx)
	chains := utils2.ContextChains(ctx)
	log.FromContext(ctx).Debugf("Running batchDeployContract()...")

	nAcc := rand.Intn(utils2.NAccounts)
	idempotency := utils.RandomString(17)
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
		return err
	}

	err = t.Load("tx.decoded", timeout)
	if err != nil {
		return err
	}

	log.FromContext(ctx).Debugf("Done: Envelope was found in tx-decoded")
	return nil
}
