package units

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"
	orchestrateclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	utils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/service/stress/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/utils/chanregistry"
)

func BatchDeployContractTest(ctx context.Context, cfg *WorkloadConfig, client orchestrateclient.OrchestrateClient, chanReg *chanregistry.ChanRegistry) error {
	logger := log.WithContext(ctx).SetComponent("deployContractTest")
	nAccount := utils.RandInt(len(cfg.accounts))
	nArtifact := utils.RandInt(len(cfg.artifacts))
	nChain := utils.RandInt(len(cfg.chains))
	idempotency := utils.RandString(30)
	evlp := tx.NewEnvelope()
	t := utils2.NewEnvelopeTracker(chanReg, evlp, idempotency)

	req := &api.DeployContractRequest{
		ChainName: cfg.chains[nChain].Name,
		Params: api.DeployContractParams{
			From:         cfg.accounts[nAccount],
			ContractName: cfg.artifacts[nArtifact],
			Args:         constructorArgs(cfg.artifacts[nArtifact]),
		},
		Labels: map[string]string{
			"id": idempotency,
		},
	}
	sReq, _ := json.Marshal(req)

	logger = logger.WithField("chain", req.ChainName).WithField("idem", idempotency)
	_, err := client.SendDeployTransaction(ctx, req)

	if err != nil {
		if !errors.IsConnectionError(err) {
			logger = logger.WithField("req", string(sReq))
		}
		logger.WithError(err).Error("failed to send transaction")
		return err
	}

	err = utils2.WaitForEnvelope(t, cfg.waitForEnvelopeTimeout)
	if err != nil {
		if !errors.IsConnectionError(err) {
			logger = logger.WithField("req", string(sReq))
		}
		logger.WithError(err).Error("envelope was not found in tx-decoded")
		return err
	}

	logger.Debug("envelope was found in tx-decoded")
	return nil
}
