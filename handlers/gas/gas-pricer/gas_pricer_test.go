package gaspricer

import (
	"context"
	"math/big"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/ethereum"
)

type MockGasPricer struct {
	t *testing.T
}

func (e *MockGasPricer) SuggestGasPrice(ctx context.Context, endpoint string) (*big.Int, error) {
	if endpoint == "error" {
		return big.NewInt(0), errors.ConnectionError("could not estimate gas")
	}
	return big.NewInt(10), nil
}

func makeGasPricerContext(i int) *engine.TxContext {
	txctx := engine.NewTxContext()
	txctx.Reset()
	txctx.Logger = log.NewEntry(log.StandardLogger())
	txctx.Envelope.Tx = &ethereum.Transaction{TxData: &ethereum.TxData{}}

	switch i % 3 {
	case 0:
		txctx.WithContext(proxy.With(txctx.Context(), "error"))
		txctx.Set("errors", 1)
		txctx.Set("result", big.NewInt(0))
	case 1:
		txctx.WithContext(proxy.With(txctx.Context(), "testURL"))
		txctx.Set("errors", 0)
		txctx.Set("result", big.NewInt(10))
	case 2:
		gp := big.NewInt(10)
		txctx.Envelope.GetTx().GetTxData().SetGasPrice(gp)
		txctx.Set("errors", 0)
		txctx.Set("result", gp)
	}
	return txctx
}

type PricerTestSuite struct {
	testutils.HandlerTestSuite
}

func (s *PricerTestSuite) SetupSuite() {
	s.Handler = Pricer(&MockGasPricer{t: s.T()})
}

func (s *PricerTestSuite) TestEstimator() {
	rounds := 100
	var txctxs []*engine.TxContext
	for i := 0; i < rounds; i++ {
		txctxs = append(txctxs, makeGasPricerContext(i))
	}

	// Handle contexts
	s.Handle(txctxs)

	for _, txctx := range txctxs {
		assert.Len(s.T(), txctx.Envelope.Errors, txctx.Get("errors").(int), "Expected right count of errors")
		for _, err := range txctx.Envelope.Errors {
			assert.Equal(s.T(), "handler.gas-pricer", err.GetComponent(), "Error should  component should have been set")
			assert.True(s.T(), errors.IsConnectionError(err), "Error should  be correct")
		}
		assert.Equal(s.T(), txctx.Get("result").(*big.Int), txctx.Envelope.GetTx().GetTxData().GetGasPriceBig(), "Expected correct Gas price")
	}
}

func TestPricer(t *testing.T) {
	suite.Run(t, new(PricerTestSuite))
}
