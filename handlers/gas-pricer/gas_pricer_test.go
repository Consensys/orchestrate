package gaspricer

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
)

type MockGasPricer struct {
	t *testing.T
}

func (e *MockGasPricer) SuggestGasPrice(ctx context.Context, chainID *big.Int) (*big.Int, error) {
	if chainID.Text(10) == "0" {
		return big.NewInt(0), fmt.Errorf("could not estimate gas")
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
		ctx.Envelope.Chain = (&common.Chain{}).SetID(big.NewInt(0))
		ctx.Set("errors", 1)
		ctx.Set("result", "")
	case 1:
		ctx.Envelope.Chain = (&common.Chain{}).SetID(big.NewInt(1))
		ctx.Set("errors", 0)
		ctx.Set("result", "0xa")
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
		assert.Equal(s.T(), txctx.Get("result").(string), txctx.Envelope.GetTx().GetTxData().GetGasPrice(), "Expected correct Gas price")
	}
}

func TestPricer(t *testing.T) {
	suite.Run(t, new(PricerTestSuite))
}
