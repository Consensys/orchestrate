package opentracing

import (
	"testing"

	oTMocktracer "github.com/opentracing/opentracing-go/mocktracer"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/tracing/opentracing"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/tracing/opentracing/mocktracer"
)

var (
	OpenTracingRootName = "Root Operation"
	OpenTracingName     = "Transaction Operation"
	MockTracer          = mocktracer.New()
	Mod                 = 5
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
		// The handler should override it
		mockSpan := MockTracer.SpanBuilder(OpenTracingRootName).Build()
		mockSpan.AttachTo(txctx)
		return txctx
	case 2:
		// There is a previous span in the Envelope Carrier
		// The handler should create a Child span of the previously existing one
		mockSpanContext := oTMocktracer.MockSpanContext{
			TraceID: 10,
			SpanID:  11,
			Sampled: true,
			Baggage: nil,
		}
		_ = MockTracer.InjectInCarrier(mockSpanContext, txctx)
		return txctx
	}
}

type TracerTestSuite struct {
	testutils.HandlerTestSuite
}

func (s *TracerTestSuite) SetupSuite() {
	opentracing.SetGlobalTracer(MockTracer)
	s.Handler = TxSpanFromBroker(MockTracer, OpenTracingName)
}

func (s *TracerTestSuite) TestTxSpanFromBroker() {

	rounds := 2
	var txctxSlice []*engine.TxContext
	for i := 0; i < rounds; i++ {
		txctxSlice = append(txctxSlice, makeTracerContext(i))
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
			assert.Equal(s.T(), OpenTracingName, span.Span.(*oTMocktracer.MockSpan).OperationName, "Expected right operationName")
		case 1:
			assert.Equal(s.T(), OpenTracingName, span.Span.(*oTMocktracer.MockSpan).OperationName, "Expected right operationName")
		case 2:
			assert.Equal(s.T(), OpenTracingName, span.Span.(*oTMocktracer.MockSpan).OperationName, "Expected right operationName")
			assert.Equal(s.T(), 44, span.Span.(*oTMocktracer.MockSpan).ParentID, "Expected right ParentID from txctx.Context")
			assert.Equal(s.T(), 43, spanContext.(oTMocktracer.MockSpanContext).TraceID, "Expected right TraceID from txctx.Envelope.Metadata")
		}
	}
}

func TestTxSpanFromBroker(t *testing.T) {
	suite.Run(t, new(TracerTestSuite))
}
