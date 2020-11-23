// +build unit

package trainjector

import (
	"testing"

	oTMocktracer "github.com/opentracing/opentracing-go/mocktracer"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/engine/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/tracing/opentracing"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/tracing/opentracing/mock"
)

var (
	OpenTracingRootName = "Root Operation"
	OpenTracingName     = "Transaction Operation"
	MockTracer          = mock.NewTracer()
	Mod                 = 2
)

func makeTracerContext(i int) *engine.TxContext {
	// Initialize context
	txctx := engine.NewTxContext().Prepare(log.NewEntry(log.StandardLogger()), nil)

	switch i % Mod {
	default:
		// There is no previous span and operation Name is not changed
		return txctx
	case 1:
		// There is a previous span in the txctx context
		// The handler should create a Child span of the previously existing one
		mockSpan := MockTracer.SpanBuilder(OpenTracingRootName).Build()
		mockSpan.AttachTo(txctx)
		return txctx
	}
}

type TracerTestSuite struct {
	testutils.HandlerTestSuite
}

func (s *TracerTestSuite) SetupSuite() {
	opentracing.SetGlobalTracer(MockTracer)
	s.Handler = TraceInjector(MockTracer, OpenTracingName)
}

func (s *TracerTestSuite) TestTxSpanFromBroker() {

	rounds := 2
	var txctxSlice []*engine.TxContext
	for i := 0; i < rounds; i++ {
		txctx := makeTracerContext(i)
		txctxSlice = append(txctxSlice, txctx)
	}

	log.Infof("TestTxSpanFromBroker: txctxSlice: %v", txctxSlice)

	// Handle contexts
	s.Handle(txctxSlice)
	for i := 0; i < rounds; i++ {
		span := opentracing.SpanFromContext(txctxSlice[i])
		spanContext, _ := MockTracer.SpanCtxFromCarrier(txctxSlice[i])

		txctxSlice[i].Logger.Debugf("TestTxSpanFromBroker: span: %v", span)
		txctxSlice[i].Logger.Debugf("TestTxSpanFromBroker: spanContext: %v", spanContext)

		switch i % Mod {
		default:
			// We just need to know if there was a panic in there
		case 1:
			assert.Equal(s.T(), 43, spanContext.(oTMocktracer.MockSpanContext).TraceID, "Expected right TraceID from txctx.Envelope.Metadata")
		}
	}
}

func TestTxSpanFromBroker(t *testing.T) {
	suite.Run(t, new(TracerTestSuite))
}
