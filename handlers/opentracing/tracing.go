package opentracing

import (
	"errors"
	"time"

	"github.com/opentracing/opentracing-go"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
)

// Errors which may occur at operation time.
var (
	ErrSpanNotFound      = errors.New("span was not found in context")
)

// TxSpanFromBroker create a new span with the given operation name and options. If a span
// is found in the TxContext and in the go Context, it will be used as the parent of the resulting span.
func TxSpanFromBroker(tracer opentracing.Tracer, defaultOperationName string) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// opts is a list of StartSpanOptions to setup the Span on creation
		opts := make([]opentracing.StartSpanOption, 0)

		// Prevent updating global operationName variable
		_operationName := defaultOperationName

		// startTime will be used to setup the Start Time of the span when created
		startTime := time.Now()

		txctx.Next()

		// find Span in TxContext.Envelope metadata, this section has been moved after the txctx.Next()
		// to be as generalistic as possible
		if spanContext, err := tracer.Extract(opentracing.TextMap, txctx.Envelope.Carrier()); err == nil {
			opts = append(opts, opentracing.FollowsFrom(spanContext))
			txctx.Logger.Tracef("TxSpanFromBroker: Spancontext in Envelope: %v", spanContext)
		} else {
			txctx.Logger.Tracef("TxSpanFromBroker: No span found during span Extraction: %v", err)
		}

		// find span context in opentracing library
		if spanParent := opentracing.SpanFromContext(txctx.Context()); spanParent != nil {
			opts = append(opts, opentracing.FollowsFrom(spanParent.Context()))
			txctx.Logger.Tracef("TxSpanFromBroker: Spanparent in Envelope: %v", spanParent)
		} else {
			txctx.Logger.Tracef("TxSpanFromBroker: No span found during span Extraction from context: %v", spanParent)
		}

		// Update span operationName if it has been created by the other middelwares
		if value, ok := txctx.Get("operationName").(string); ok {
			_operationName = value
		}

		// Add in StartSpanOptions the starting time previously set
		opts = append(opts, opentracing.StartTime(startTime))

		span := tracer.StartSpan(_operationName, opts...)
		defer span.Finish()

		txctx.WithContext(opentracing.ContextWithSpan(txctx.Context(), span))

		if err := tracer.Inject(span.Context(), opentracing.TextMap, txctx.Envelope.Carrier()); err != nil {
			txctx.Logger.Errorf("TxSpanFromBroker: Error during span Injection %v", err)
		}
	}
}
