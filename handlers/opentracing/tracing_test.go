package opentracing

import (
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine/testutils"
)

var (
	OpenTracingRootName = "Root Operation"
	OpenTracingName     = "Transaction Operation"
	MockTracer          = mocktracer.New()
	Mod                 = 3 // TODO : IMPORTANT : OpenTracing is NOT thread Safe, modify Mod = 5 to apply all tests
)

func makeTracerContext(i int) *engine.TxContext {
	// Initialize context
	txctx := engine.NewTxContext().Prepare([]engine.HandlerFunc{}, log.NewEntry(log.StandardLogger()), nil)

	txctx.Reset()
	switch i % Mod {
	default:
		return txctx
	case 1:
		txctx.Set("operationName", "I love Crafting")
		return txctx
	case 2:
		mockSpan := MockTracer.StartSpan(OpenTracingRootName)
		txctx.WithContext(opentracing.ContextWithSpan(txctx.Context(), mockSpan))
		return txctx
	case 3:
		mockSpanContext := mocktracer.MockSpanContext{
			TraceID: 10,
			SpanID:  11,
			Sampled: true,
			Baggage: nil,
		}

		MockTracer.Inject(mockSpanContext, opentracing.TextMap, txctx.Envelope.Carrier())

		return txctx
	case 4:
		mockSpanContext := mocktracer.MockSpanContext{
			TraceID: 10,
			SpanID:  11,
			Sampled: true,
			Baggage: nil,
		}

		MockTracer.Inject(mockSpanContext, opentracing.TextMap, txctx.Envelope.Carrier())

		mockSpan := MockTracer.StartSpan(OpenTracingRootName, opentracing.FollowsFrom(mockSpanContext))
		txctx.WithContext(opentracing.ContextWithSpan(txctx.Context(), mockSpan))

		return txctx
	}
}

type TracerTestSuite struct {
	testutils.HandlerTestSuite
}

func (suite *TracerTestSuite) SetupSuite() {
	opentracing.SetGlobalTracer(MockTracer)
	suite.Handler = TxSpanFromBroker(MockTracer, OpenTracingName)
}

func (suite *TracerTestSuite) TestTxSpanFromBroker() {

	rounds := 1
	var txctxSlice []*engine.TxContext
	for i := 0; i < rounds; i++ {
		txctxSlice = append(txctxSlice, makeTracerContext(i))
	}

	// Handle contexts
	suite.Handle(txctxSlice)
	for i := 0; i < rounds; i++ {
		span := opentracing.SpanFromContext(txctxSlice[i].Context())
		spanContext, _ := MockTracer.Extract(opentracing.TextMap, txctxSlice[i].Envelope.Carrier())

		switch i % Mod {
		default:
			assert.Equal(suite.T(), GenericOperationName, span.(*mocktracer.MockSpan).OperationName, "Expected right operationName")
			break
		case 1:
			assert.Equal(suite.T(), "I love Crafting", span.(*mocktracer.MockSpan).OperationName, "Expected right operationName")
			break
		case 2:
			assert.Equal(suite.T(), GenericOperationName, span.(*mocktracer.MockSpan).OperationName, "Expected right operationName")
			assert.Equal(suite.T(), 44, span.(*mocktracer.MockSpan).ParentID, "Expected right ParentID from txctx.Context")
			assert.Equal(suite.T(), 43, spanContext.(mocktracer.MockSpanContext).TraceID, "Expected right TraceID from txctx.Envelope.Metadata")
			assert.Equal(suite.T(), 46, spanContext.(mocktracer.MockSpanContext).SpanID, "Expected right SpanID from txctx.Envelope.Metadata")
			break
		case 3, 4: // TODO : IMPORTANT : OpenTracing is NOT thread Safe, modify Mod = 5 to apply all tests
			assert.Equal(suite.T(), GenericOperationName, span.(*mocktracer.MockSpan).OperationName, "Expected right operationName")
			assert.Equal(suite.T(), 11, span.(*mocktracer.MockSpan).ParentID, "Expected right operationName")
			assert.Equal(suite.T(), 10, spanContext.(mocktracer.MockSpanContext).TraceID, "Expected right TraceID from txctx.Envelope.Metadata")
			assert.Equal(suite.T(), 46, spanContext.(mocktracer.MockSpanContext).SpanID, "Expected right SpanID from txctx.Envelope.Metadata")
			break
		}
	}
}

func TestTxSpanFromBroker(t *testing.T) {
	suite.Run(t, new(TracerTestSuite))
}
