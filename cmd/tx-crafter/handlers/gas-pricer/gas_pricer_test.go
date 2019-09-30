package gaspricer

import (
	"context"
	"math/big"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/types/chain"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/types/ethereum"
)

type MockGasPricer struct {
	t *testing.T
}

func (e *MockGasPricer) SuggestGasPrice(ctx context.Context, chainID *big.Int) (*big.Int, error) {
	if chainID.Text(10) == "0" {
		return big.NewInt(0), errors.ConnectionError("could not estimate gas")
	}
	return big.NewInt(10), nil
}

func makeGasPricerContext(i int) *engine.TxContext {
	ctx := engine.NewTxContext()
	ctx.Reset()
	ctx.Logger = log.NewEntry(log.StandardLogger())
	ctx.Envelope.Tx = &ethereum.Transaction{TxData: &ethereum.TxData{}}

	switch i % 2 {
	case 0:
		ctx.Envelope.Chain = (&chain.Chain{}).SetID(big.NewInt(0))
		ctx.Set("errors", 1)
		ctx.Set("result", big.NewInt(0))
	case 1:
		ctx.Envelope.Chain = (&chain.Chain{}).SetID(big.NewInt(1))
		ctx.Set("errors", 0)
		ctx.Set("result", big.NewInt(10))
	}
	return ctx
}

type PricerTestSuite struct {
	testutils.HandlerTestSuite
}

func (s *PricerTestSuite) SetupSuite() {
	s.Handler = Pricer(&MockGasPricer{t: s.T()})
}

func (s *PricerTestSuite) TestEstimator() {
	rounds := 100
	txctxs := []*engine.TxContext{}
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
		assert.Equal(s.T(), txctx.Get("result").(*big.Int), txctx.Envelope.GetTx().GetTxData().GetGasPrice().Value(), "Expected correct Gas price")
	}
}

func TestPricer(t *testing.T) {
	suite.Run(t, new(PricerTestSuite))
}
