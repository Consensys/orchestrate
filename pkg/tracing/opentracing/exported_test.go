package opentracing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber/jaeger-client-go"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	tracer := GetGlobalTracer()
	assert.NotNil(t, tracer, "Jaeger Tracer should not be nil")

	_, ok := tracer.Internal.(*jaeger.Tracer)
	assert.True(t, ok, "Jaeger Tracer should cast to jaeger.Tracer")
}
