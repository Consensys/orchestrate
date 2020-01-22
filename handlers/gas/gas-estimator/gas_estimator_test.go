package gasestimator

import (
	"context"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"

	eth "github.com/ethereum/go-ethereum"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/ethereum"
)

type MockGasEstimator struct {
	t *testing.T
}

func (e *MockGasEstimator) EstimateGas(ctx context.Context, endpoint string, call *eth.CallMsg) (uint64, error) { // nolint:gocritic
	if endpoint == "error" {
		return 0, errors.ConnectionError("could not estimate gas").SetComponent("mock")
	}
	return 18, nil
}

func makeGasEstimatorContext(i int) *engine.TxContext {
	txctx := engine.NewTxContext()
	txctx.Reset()
	txctx.Logger = log.NewEntry(log.StandardLogger())
	txctx.Envelope.From = &ethereum.Account{}
	txctx.Envelope.Tx = &ethereum.Transaction{TxData: &ethereum.TxData{}}

	switch i % 2 {
	case 0:
		txctx.WithContext(proxy.With(txctx.Context(), "error"))
		txctx.Set("errors", 1)
		txctx.Set("result", uint64(0))
	case 1:
		txctx.WithContext(proxy.With(txctx.Context(), "testURL"))
		txctx.Set("errors", 0)
		txctx.Set("result", uint64(18))
	}
	return txctx
}

type EstimatorTestSuite struct {
	testutils.HandlerTestSuite
}

func (s *EstimatorTestSuite) SetupSuite() {
	s.Handler = Estimator(&MockGasEstimator{t: s.T()})
}

func (s *EstimatorTestSuite) TestEstimator() {
	rounds := 100
	var txctxs []*engine.TxContext
	for i := 0; i < rounds; i++ {
		txctxs = append(txctxs, makeGasEstimatorContext(i))
	}

	// Handle contexts
	s.Handle(txctxs)

	for _, txctx := range txctxs {
		assert.Len(s.T(), txctx.Envelope.Errors, txctx.Get("errors").(int), "Expected right count of errors")
		for _, err := range txctx.Envelope.Errors {
			assert.Equal(s.T(), "handler.gas-estimator.mock", err.GetComponent(), "Error  component should have been set")
			assert.True(s.T(), errors.IsConnectionError(err), "Error should  be correct")
		}
		assert.Equal(s.T(), txctx.Get("result").(uint64), txctx.Envelope.GetTx().GetTxData().GetGas(), "Expected correct payload")
	}
}

func TestEstimator(t *testing.T) {
	suite.Run(t, new(EstimatorTestSuite))
}
