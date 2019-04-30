package jaeger

import (
	"context"
	"sync"

	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go/config"
)

var (
	cfg      *config.Configuration
	initOnce = &sync.Once{}
)

// Init initialize tracer
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if cfg == nil {
			cfg = NewConfig()
		}

		tracer, _, err := cfg.NewTracer()
		if err != nil {
			log.WithError(err).Fatal("opentracing: could initialize jaeger tracer")
		}

		// Set Open tracing global tracer
		opentracing.SetGlobalTracer(tracer)
	})
}
