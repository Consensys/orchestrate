package jaeger

import (
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
)

// InitTracer initialize tracer
func InitTracer() (opentracing.Tracer, io.Closer) {
	cfg := NewConfig()
	tracer, closer, err := cfg.NewTracer()

	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Tracer Jaeger: %v\n", err))
	}

	opentracing.SetGlobalTracer(tracer)

	return tracer, closer
}
