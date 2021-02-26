package units

import (
	"context"

	"github.com/ConsenSys/orchestrate/pkg/encoding/json"
	"github.com/ConsenSys/orchestrate/pkg/errors"
	orchestrateclient "github.com/ConsenSys/orchestrate/pkg/sdk/client"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
	"github.com/ConsenSys/orchestrate/pkg/types/api"
	"github.com/ConsenSys/orchestrate/pkg/types/tx"
	"github.com/ConsenSys/orchestrate/pkg/utils"
	utils2 "github.com/ConsenSys/orchestrate/tests/service/stress/utils"
	utils3 "github.com/ConsenSys/orchestrate/tests/utils"
	"github.com/ConsenSys/orchestrate/tests/utils/chanregistry"
)

func BatchDeployContractTest(ctx context.Context, cfg *WorkloadConfig, client orchestrateclient.OrchestrateClient, chanReg *chanregistry.ChanRegistry) error {
	logger := log.WithContext(ctx).SetComponent("stress-test.deploy-contract")
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
		logger.WithField("topic", utils3.TxDecodedTopicKey).WithError(err).Error("envelope was not found in topic")
		return err
	}

	logger.WithField("topic", utils3.TxDecodedTopicKey).Debug("envelope was found in topic")
	return nil
}
