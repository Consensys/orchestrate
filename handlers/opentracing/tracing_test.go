package opentracing

import (
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine/testutils"
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
		// There is no previous span but operation Name is changed
		txctx.Set("operationName", "I love Crafting")
		return txctx
	case 2:
		// There is a previous span in the txctx context
		// The handler should create a Child span of the previously existing one
		mockSpan := MockTracer.StartSpan(OpenTracingRootName)
		txctx.WithContext(opentracing.ContextWithSpan(txctx.Context(), mockSpan))
		return txctx
	case 3:
		// There is a previous span in the Envelope Carrier
		// The handler should create a Child span of the previously existing one
		mockSpanContext := mocktracer.MockSpanContext{
			TraceID: 10,
			SpanID:  11,
			Sampled: true,
			Baggage: nil,
		}

		_ = MockTracer.Inject(mockSpanContext, opentracing.TextMap, txctx.Envelope.Carrier())

		return txctx
	case 4:
		// There are two previous spans, one in the txctx context and one in the Envelope Carrier
		// The handler should create a Child span of the span in the Envelope
		mockSpanContext := mocktracer.MockSpanContext{
			TraceID: 10,
			SpanID:  11,
			Sampled: true,
			Baggage: nil,
		}

		_ = MockTracer.Inject(mockSpanContext, opentracing.TextMap, txctx.Envelope.Carrier())

		mockSpan := MockTracer.StartSpan(OpenTracingRootName, opentracing.FollowsFrom(mockSpanContext))
		txctx.WithContext(opentracing.ContextWithSpan(txctx.Context(), mockSpan))

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

	rounds := 5
	txctxSlice := []*engine.TxContext{}
	for i := 0; i < rounds; i++ {
		txctxSlice = append(txctxSlice, makeTracerContext(i))
	}

	log.Infof("TestTxSpanFromBroker: txctxSlice: %v", txctxSlice)

	// Handle contexts
	s.Handle(txctxSlice)
	for i := 0; i < rounds; i++ {
		span := opentracing.SpanFromContext(txctxSlice[i].Context())
		spanContext, _ := MockTracer.Extract(opentracing.TextMap, txctxSlice[i].Envelope.Carrier())

		txctxSlice[i].Logger.Debugf("TestTxSpanFromBroker: span: %v", span)
		txctxSlice[i].Logger.Debugf("TestTxSpanFromBroker: spanContext: %v", spanContext)

		switch i % Mod {
		default:
			assert.Equal(s.T(), OpenTracingName, span.(*mocktracer.MockSpan).OperationName, "Expected right operationName")
		case 1:
			assert.Equal(s.T(), "I love Crafting", span.(*mocktracer.MockSpan).OperationName, "Expected right operationName")
		case 2:
			assert.Equal(s.T(), OpenTracingName, span.(*mocktracer.MockSpan).OperationName, "Expected right operationName")
			assert.Equal(s.T(), 44, span.(*mocktracer.MockSpan).ParentID, "Expected right ParentID from txctx.Context")
			assert.Equal(s.T(), 43, spanContext.(mocktracer.MockSpanContext).TraceID, "Expected right TraceID from txctx.Envelope.Metadata")
		case 3:
			assert.Equal(s.T(), OpenTracingName, span.(*mocktracer.MockSpan).OperationName, "Expected right operationName")
			assert.Equal(s.T(), 11, span.(*mocktracer.MockSpan).ParentID, "Expected right operationName")
			assert.Equal(s.T(), 10, spanContext.(mocktracer.MockSpanContext).TraceID, "Expected right TraceID from txctx.Envelope.Metadata")
		case 4:
			assert.Equal(s.T(), OpenTracingName, span.(*mocktracer.MockSpan).OperationName, "Expected right operationName")
			assert.Equal(s.T(), 11, span.(*mocktracer.MockSpan).ParentID, "Expected right operationName")
			assert.Equal(s.T(), 10, spanContext.(mocktracer.MockSpanContext).TraceID, "Expected right TraceID from txctx.Envelope.Metadata")
		}
	}
}

func TestTxSpanFromBroker(t *testing.T) {
	suite.Run(t, new(TracerTestSuite))
}
