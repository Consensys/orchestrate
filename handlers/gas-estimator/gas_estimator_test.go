package gasestimator

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	eth "github.com/ethereum/go-ethereum"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/chain"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
)

type MockGasEstimator struct {
	t *testing.T
}

func (e *MockGasEstimator) EstimateGas(ctx context.Context, chainID *big.Int, call *eth.CallMsg) (uint64, error) { // nolint:gocritic
	if chainID.Text(10) == "0" {
		return 0, fmt.Errorf("could not estimate gas")
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
		txctx.Envelope.Chain = (&chain.Chain{}).SetID(big.NewInt(0))
		txctx.Set("errors", 1)
		txctx.Set("result", uint64(0))
	case 1:
		txctx.Envelope.Chain = (&chain.Chain{}).SetID(big.NewInt(1))
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
	txctxs := []*engine.TxContext{}
	for i := 0; i < rounds; i++ {
		txctxs = append(txctxs, makeGasEstimatorContext(i))
	}

	// Handle contexts
	s.Handle(txctxs)

	for _, txctx := range txctxs {
		assert.Len(s.T(), txctx.Envelope.Errors, txctx.Get("errors").(int), "Expected right count of errors")
		assert.Equal(s.T(), txctx.Get("result").(uint64), txctx.Envelope.GetTx().GetTxData().GetGas(), "Expected correct payload")
	}
}

func TestEstimator(t *testing.T) {
	suite.Run(t, new(EstimatorTestSuite))
}
