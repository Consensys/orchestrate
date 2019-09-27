package logger

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine/testutils"
)

func makeLoggerContext() *engine.TxContext {
	txctx := engine.NewTxContext().Prepare(log.NewEntry(log.StandardLogger()), nil)
	txctx.Set("errors", 0)
	return txctx
}

type LoggerTestSuite struct {
	testutils.HandlerTestSuite
}

func (s *LoggerTestSuite) SetupSuite() {
	s.Handler = Logger
}

func (s *LoggerTestSuite) TestLogger() {
	rounds := 100
	txctxs := []*engine.TxContext{}
	for i := 0; i < rounds; i++ {
		txctxs = append(txctxs, makeLoggerContext())
	}

	// Handle contexts
	s.Handle(txctxs)

	for _, txctx := range txctxs {
		assert.Len(s.T(), txctx.Envelope.Errors, txctx.Get("errors").(int), "Expected right count of errors")
	}
}

func TestLogger(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}
