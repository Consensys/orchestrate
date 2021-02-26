package units

import (
	"context"

	"github.com/ConsenSys/orchestrate/pkg/encoding/json"
	"github.com/ConsenSys/orchestrate/pkg/errors"
	orchestrateclient "github.com/ConsenSys/orchestrate/pkg/sdk/client"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
	"github.com/ConsenSys/orchestrate/pkg/types/api"
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/ConsenSys/orchestrate/pkg/types/tx"
	"github.com/ConsenSys/orchestrate/pkg/utils"
	"github.com/ConsenSys/orchestrate/tests/service/stress/assets"
	utils2 "github.com/ConsenSys/orchestrate/tests/service/stress/utils"
	utils3 "github.com/ConsenSys/orchestrate/tests/utils"
	"github.com/ConsenSys/orchestrate/tests/utils/chanregistry"
)

func BatchPrivateTxsTest(ctx context.Context, cfg *WorkloadConfig, client orchestrateclient.OrchestrateClient, chanReg *chanregistry.ChanRegistry) error {
	logger := log.WithContext(ctx).SetComponent("stress-test.private-txs")

	account := cfg.accounts[utils.RandInt(len(cfg.accounts))]
	contractName := cfg.artifacts[utils.RandInt(len(cfg.artifacts))]
	chain := cfg.chains[utils.RandInt(len(cfg.chains))]
	privacyGroup := cfg.privacyGroups[utils.RandInt(len(cfg.privacyGroups))]
	privateFrom := chain.PrivNodeAddress[utils.RandInt(len(chain.PrivNodeAddress))]
	idempotency := utils.RandString(30)

	evlp := tx.NewEnvelope()
	t := utils2.NewEnvelopeTracker(chanReg, evlp, idempotency)

	req := &api.DeployContractRequest{
		ChainName: chain.Name,
		Params: api.DeployContractParams{
			From:         account,
			ContractName: contractName,
			Args:         constructorArgs(contractName),
			PrivateFrom:  privateFrom,
			Protocol:     entities.OrionChainType,
		},
		Labels: map[string]string{
			"id": idempotency,
		},
	}

	usePrivacyGroup := canUsePrivacyGroup(chain.PrivNodeAddress, &privacyGroup)
	if usePrivacyGroup {
		req.Params.PrivacyGroupID = privacyGroup.ID
	} else {
		size := len(privacyGroup.Nodes)
		req.Params.PrivateFor = privacyGroup.Nodes[0 : size-1]
	}

	sReq, _ := json.Marshal(req)
	logger = logger.WithField("chain", req.ChainName).WithField("idem", idempotency)
	logger.Debug("sending private tx to deploy contract")

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
		logger.WithError(err).Error("failed to fetch envelope")
		return err
	}

	logger.WithField("topic", utils3.TxDecodedTopicKey).Debug("envelope was found")
	return nil
}

func canUsePrivacyGroup(chainPrivNodes []string, pGroup *assets.PrivacyGroup) bool {
	for _, cAddr := range chainPrivNodes {
		for _, gAddr := range pGroup.Nodes {
			if cAddr == gAddr {
				return true
			}
		}
	}

	return false
}
