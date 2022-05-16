package units

import (
	"context"
	"encoding/json"

	"github.com/consensys/orchestrate/pkg/errors"
	orchestrateclient "github.com/consensys/orchestrate/pkg/sdk/client"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/pkg/types/tx"
	"github.com/consensys/orchestrate/pkg/utils"
	utils2 "github.com/consensys/orchestrate/tests/service/stress/utils"
	utils3 "github.com/consensys/orchestrate/tests/utils"
	"github.com/consensys/orchestrate/tests/utils/chanregistry"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

func SendContractTxsTest(ctx context.Context, cfg *WorkloadConfig, client orchestrateclient.OrchestrateClient, chanReg *chanregistry.ChanRegistry) error {
	logger := log.WithContext(ctx).SetComponent("stress-test.send-contract-txs")
	nAccount := utils.RandInt(len(cfg.accounts))
	nChain := utils.RandInt(len(cfg.chains))
	idempotency := utils.RandString(30)
	evlp := tx.NewEnvelope()
	t := utils2.NewEnvelopeTracker(chanReg, evlp, idempotency)

	// @TODO Read values from configuration or from context
	toAddr := ethcommon.HexToAddress("0xFf80849F797a5feBC96F1737dc78135a79DaF83E")
	req := &api.SendTransactionRequest{
		ChainName: cfg.chains[nChain].Name,
		Params: api.TransactionParams{
			From:            &cfg.accounts[nAccount],
			To:              &toAddr,
			ContractName:    "Counter",
			MethodSignature: "increment(uint256)",
			Args:            []interface{}{utils.RandInt(100)},
		},
		Labels: map[string]string{
			"id": idempotency,
		},
	}
	sReq, _ := json.Marshal(req)

	logger = logger.WithField("chain", req.ChainName).WithField("idem", idempotency)
	_, err := client.SendContractTransaction(ctx, req)

	if err != nil {
		if !errors.IsConnectionError(err) {
			logger = logger.WithField("req", string(sReq))
		}
		logger.WithError(err).Error("failed to send contract transaction")
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
