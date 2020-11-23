package mock

import (
	"github.com/opentracing/opentracing-go/mocktracer"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/tracing/opentracing"
)

// NewTracer returns a mocked tracer instance
func NewTracer() *opentracing.Tracer {
	return &opentracing.Tracer{Internal: mocktracer.New()}
}
