// +build unit

package opentracing

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber/jaeger-client-go"
)

func TestInit(t *testing.T) {
	_ = os.Setenv("JAEGER_ENABLED", "true")
	Init(context.Background())
	_ = os.Unsetenv("JAEGER_ENABLED")
	
	tracer := GetGlobalTracer()
	assert.NotNil(t, tracer, "Jaeger Tracer should not be nil")

	_, ok := tracer.Internal.(*jaeger.Tracer)
	assert.True(t, ok, "Jaeger Tracer should cast to jaeger.Tracer")
}
