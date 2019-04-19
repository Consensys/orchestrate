package opentracing

import (
	"errors"

	"github.com/opentracing/opentracing-go"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"

	log "github.com/sirupsen/logrus"
)

// Errors which may occur at operation time.
var (
	ErrSpanNotFound      = errors.New("span was not found in context")
	GenericOperationName = "Transaction Operation"
)

// TxSpanFromBroker create a new span with the given operation name and options. If a span
// is found in the TxContext and in the go Context, it will be used as the parent of the resulting span.
func TxSpanFromBroker(tracer opentracing.Tracer, operationName string, opts ...opentracing.StartSpanOption) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// find Span in TxContext.Envelope metadata
		if spanContext, err := tracer.Extract(opentracing.TextMap, txctx.Envelope.Carrier()); err == nil {
			opts = append(opts, opentracing.ChildOf(spanContext))
		}

		// find span context in opentracing library
		if spanParent := opentracing.SpanFromContext(txctx.Context()); spanParent != nil {
			opts = append(opts, opentracing.ChildOf(spanParent.Context()))
		}

		span := tracer.StartSpan(operationName, opts...)
		defer span.Finish()

		txctx.WithContext(opentracing.ContextWithSpan(txctx.Context(), span))

		txctx.Next()

		if value, ok := txctx.Get("operationName").(string); ok {
			span.SetOperationName(value)
			txctx.WithContext(opentracing.ContextWithSpan(txctx.Context(), span))
		}

		if err := span.Tracer().Inject(span.Context(), opentracing.TextMap, txctx.Envelope.Carrier()); err != nil {
			log.Errorf("Error during span Injection %v", err)
		}
	}
}
