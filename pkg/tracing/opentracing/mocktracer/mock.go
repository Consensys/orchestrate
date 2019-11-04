package mocktracer

import (
	"github.com/opentracing/opentracing-go/mocktracer"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/tracing/opentracing"
)

// New returns a mocked tracer instance
func New() *opentracing.Tracer {
	return &opentracing.Tracer{Internal: mocktracer.New()}
}
