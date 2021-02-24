package mock

import (
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/ConsenSys/orchestrate/pkg/tracing/opentracing"
)

// NewTracer returns a mocked tracer instance
func NewTracer() *opentracing.Tracer {
	return &opentracing.Tracer{Internal: mocktracer.New()}
}
