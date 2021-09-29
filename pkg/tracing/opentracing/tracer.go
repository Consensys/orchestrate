package opentracing

import (
	"time"

	"github.com/opentracing/opentracing-go"

	"github.com/consensys/orchestrate/pkg/engine"
	"github.com/consensys/orchestrate/pkg/errors"
)

// Tracer is a facade to manage OpenTracing utilities more easily
type Tracer struct {
	Internal opentracing.Tracer
}

// NewTracer is a default constructor for a tracer
func NewTracer(internal opentracing.Tracer) *Tracer {
	return &Tracer{Internal: internal}
}

// SpanBuilder is a set of OpenTracing options
type SpanBuilder struct {
	Opts          []opentracing.StartSpanOption
	Tracer        *Tracer
	OperationName string
}

// Span is a wrapper around jaeger.Span that can interact with txctx
type Span struct {
	opentracing.Span
}

// SpanBuilder is a default constructor for
func (t *Tracer) SpanBuilder(operationName string) *SpanBuilder {
	return &SpanBuilder{
		Opts:          make([]opentracing.StartSpanOption, 0),
		OperationName: operationName,
		Tracer:        t,
	}
}

// InjectFromContext injects a trace in the carrier from context
func (t *Tracer) InjectFromContext(txctx *engine.TxContext) error {
	span := SpanFromContext(txctx)
	if span.Span == nil {
		return errors.NotFoundError("Could not find span in context")
	}
	return t.InjectInCarrier(span.Span.Context(), txctx)
}

// SpanFromContext returns a wrapped opentracing.Span object from txctx
func SpanFromContext(txctx *engine.TxContext) *Span {
	span := opentracing.SpanFromContext(txctx.Context())
	return &Span{span}
}

// SpanCtxFromCarrier returns a wrapped opentracing.Span object from txctx
func (t *Tracer) SpanCtxFromCarrier(txctx *engine.TxContext) (opentracing.SpanContext, error) {
	return t.Internal.Extract(
		opentracing.TextMap,
		txctx.Envelope.Carrier())
}

// InjectInCarrier a span context in the txctx carrier
func (t *Tracer) InjectInCarrier(spanCtx opentracing.SpanContext, txctx *engine.TxContext) error {
	return t.Internal.Inject(
		spanCtx,
		opentracing.TextMap,
		txctx.Envelope.Carrier(),
	)
}

// FollowingFromCarrier extracts a SpanContext from the Carrier if it exists
func (s *SpanBuilder) FollowingFromCarrier(txctx *engine.TxContext) *SpanBuilder {
	spanContext, err := s.Tracer.SpanCtxFromCarrier(txctx)

	if err != nil || spanContext == nil {
		// Ignore and do nothing
		return s
	}

	s.Opts = append(s.Opts, opentracing.FollowsFrom(spanContext))
	return s
}

// FollowingFromContext extracts a SpanContext from txctx and returns an option
func (s *SpanBuilder) FollowingFromContext(txctx *engine.TxContext) *SpanBuilder {
	s.Opts = append(s.Opts, opentracing.FollowsFrom(
		SpanFromContext(txctx).Context(),
	))
	return s
}

// StartingAt attaches a starting time for the created span
func (s *SpanBuilder) StartingAt(t time.Time) *SpanBuilder {
	s.Opts = append(s.Opts, opentracing.StartTime(t))
	return s
}

// Build the span with the previously provided options into txctx
func (s *SpanBuilder) Build() *Span {
	span := s.Tracer.Internal.StartSpan(
		s.OperationName,
		s.Opts...)

	return &Span{span}
}

// AttachTo attaches a span to txctx
func (s *Span) AttachTo(txctx *engine.TxContext) {
	txctx.WithContext(opentracing.ContextWithSpan(txctx.Context(), s.Span))
}

// Finish the current span
func (s *Span) Finish() {
	s.Span.Finish()
}
