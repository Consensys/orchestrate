package faucet

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	faucettypes "gitlab.com/ConsenSys/client/fr/core-stack/service/faucet.git/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
)

type MockFaucet struct {
	t *testing.T
}

func (f *MockFaucet) Credit(ctx context.Context, r *faucettypes.Request) (*big.Int, bool, error) {
	if r.ChainID.Text(10) == "0" {
		return big.NewInt(0), false, fmt.Errorf("could not credit")
	}
	return r.Amount, true, nil
}

func makeFaucetContext(i int) *engine.TxContext {
	txctx := engine.NewTxContext()
	txctx.Reset()
	txctx.Logger = log.NewEntry(log.StandardLogger())
	switch i % 2 {
	case 0:
		txctx.Envelope.Chain = &common.Chain{}
		txctx.Set("errors", 1)
	case 1:
		txctx.Envelope.Chain = (&common.Chain{}).SetID(big.NewInt(1))
		txctx.Set("errors", 0)
	}
	return txctx
}

type FaucetTestSuite struct {
	testutils.HandlerTestSuite
}

func (s *FaucetTestSuite) SetupSuite() {
	s.Handler = Faucet(&MockFaucet{t: s.T()})
}

func (s *FaucetTestSuite) TestFaucet() {
	rounds := 100
	txctxs := []*engine.TxContext{}
	for i := 0; i < rounds; i++ {
		txctxs = append(txctxs, makeFaucetContext(i))
	}

	// Handle contexts
	s.Handle(txctxs)

	for _, txctx := range txctxs {
		assert.Len(s.T(), txctx.Envelope.Errors, txctx.Get("errors").(int), "Expected right count of errors")
	}
}

func TestFaucet(t *testing.T) {
	suite.Run(t, new(FaucetTestSuite))
}
