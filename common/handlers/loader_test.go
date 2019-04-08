package handlers

import (
	"fmt"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/envelope"
)

type MockUnmarshaller struct {
	t *testing.T
}

func (u *MockUnmarshaller) Unmarshal(msg interface{}, envelope *envelope.Envelope) error {
	if msg.(string) == "error" {
		return fmt.Errorf("Could not unmarshall")
	}
	return nil
}

func makeLoaderContext(i int) *engine.TxContext {
	// Initialize context
	txctx := engine.NewTxContext().Prepare([]engine.HandlerFunc{}, log.NewEntry(log.StandardLogger()), nil)

	switch i % 2 {
	case 0:
		txctx.Msg = "error"
		txctx.Keys["errors"] = 1
	case 1:
		txctx.Msg = "valid"
		txctx.Keys["errors"] = 0
	}

	return txctx
}

type LoaderTestSuite struct {
	testutils.HandlerTestSuite
}

func (suite *LoaderTestSuite) SetupSuite() {
	suite.Handler = Loader(&MockUnmarshaller{t: suite.T()})
}

func (suite *LoaderTestSuite) TestLoader() {
	rounds := 100
	txctxs := []*engine.TxContext{}
	for i := 0; i < rounds; i++ {
		txctxs = append(txctxs, makeLoaderContext(i))
	}

	// Handle contexts
	suite.Handle(txctxs)

	for _, txctx := range txctxs {
		assert.Len(suite.T(), txctx.Envelope.Errors, txctx.Keys["errors"].(int), "Expected right count of errors")
	}
}

func TestLoader(t *testing.T) {
	suite.Run(t, new(LoaderTestSuite))
}
