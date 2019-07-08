package jaeger

import (
	"context"
	"reflect"
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
			log.WithError(err).Fatal("open-tracing: could initialize jaeger tracer")
		}

		// Get tracer specific values
		serviceName := reflect.ValueOf(tracer).Elem().FieldByName("serviceName")
		tags := reflect.ValueOf(tracer).Elem().FieldByName("tags")

		log.Infof("jaeger: agent ready for open-tracing")
		log.Infof("jaeger: service name: %v", serviceName)
		log.Infof("jaeger: tags: %v", tags)

		// Set Open tracing global tracer
		opentracing.SetGlobalTracer(tracer)
	})
}
