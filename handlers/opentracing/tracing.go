package opentracing

import (
	"errors"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/tracing/opentracing"
)

// Errors which may occur at operation time.
var (
	ErrSpanNotFound = errors.New("span was not found in context")
)

// TxSpanFromBroker create a new span with the given operation name and options. If a span
// is found in the Envelope and in the go Context, it will be used as the parent of the resulting span.
func TxSpanFromBroker(tracer *opentracing.Tracer, defaultOperationName string) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// If there is a span in the context, finish it
		prevSpan := opentracing.SpanFromContext(txctx)
		if prevSpan.Span != nil {
			prevSpan.Finish()
		}
		// Prevent updating global operationName variable
		_operationName := defaultOperationName
		// Builds a span following from a carried trace if it exists
		span := tracer.SpanBuilder(_operationName).
			FollowingFromCarrier(txctx).
			StartingAt(time.Now()).
			Build()
		// Attach the span to txctx
		span.AttachTo(txctx)
		// Run the next middlewares
		txctx.Next()
		// Close the span and report
		span.Finish()
	}
}
