package jaeger

import (
	"context"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.com/uber/jaeger-client-go"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	_, ok := opentracing.GlobalTracer().(*jaeger.Tracer)
	assert.True(t, ok, "Jaeger Tracer should have been set")
}
