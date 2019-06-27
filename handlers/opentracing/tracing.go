package opentracing

import (
	"time"
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
func TxSpanFromBroker(tracer opentracing.Tracer, operationName string) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// opts is a list of StartSpanOptions to setup the Span on creation
		opts := make([]opentracing.StartSpanOption, 0)

		// startTime will be used to setup the Start Time of the span when created
		startTime := time.Now()

		txctx.Next()

		// find Span in TxContext.Envelope metadata, this section has been moved after the txctx.Next()
		// to be as generalistic as possible
		if spanContext, err := tracer.Extract(opentracing.TextMap, txctx.Envelope.Carrier()); err == nil {
			opts = append(opts, opentracing.ChildOf(spanContext))
			log.Debugf("Spancontext in Enveloppe: %v", spanContext)
		} else {
			log.Debugf("No span found during span Extraction: %v", err)
		}

		// find span context in opentracing library
		if spanParent := opentracing.SpanFromContext(txctx.Context()); spanParent != nil {
			opts = append(opts, opentracing.ChildOf(spanParent.Context()))
			log.Debugf("Spanparent in Enveloppe: %v", spanParent)
		} else {
			log.Debugf("No span found during span Extraction from context: %v", spanParent)
		}
		
		// Update span operationName if it has been created by the other middelwares
		if value, ok := txctx.Get("operationName").(string); ok {
			operationName = value
		}

		// Add in StartSpanOptions the starting time previously set 
		opts = append(opts, opentracing.StartTime(startTime))

		span := tracer.StartSpan(operationName, opts...)
		defer span.Finish()

		// TODO: why using a different tracer between Extract and Inject
		if err := span.Tracer().Inject(span.Context(), opentracing.TextMap, txctx.Envelope.Carrier()); err != nil {
			log.Errorf("Error during span Injection %v", err)
		}
	}
}
