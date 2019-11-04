package faucet

import (
	"context"
	"math/big"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	faucettypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/chain"
)

type MockFaucet struct {
	t *testing.T
}

func (f *MockFaucet) Credit(ctx context.Context, r *faucettypes.Request) (*big.Int, bool, error) {
	if r.ChainID.Text(10) == "0" {
		return big.NewInt(0), false, errors.FaucetWarning("could not credit").SetComponent("mock")
	}
	return r.Amount, true, nil
}

func makeFaucetContext(i int) *engine.TxContext {
	txctx := engine.NewTxContext()
	txctx.Reset()
	txctx.Logger = log.NewEntry(log.StandardLogger())
	switch i % 2 {
	case 0:
		txctx.Envelope.Chain = chain.FromInt(0)
	case 1:
		txctx.Envelope.Chain = chain.FromInt(10)
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
	var txctxs []*engine.TxContext
	for i := 0; i < rounds; i++ {
		txctxs = append(txctxs, makeFaucetContext(i))
	}

	// Handle contexts
	s.Handle(txctxs)
}

func TestFaucet(t *testing.T) {
	suite.Run(t, new(FaucetTestSuite))
}
