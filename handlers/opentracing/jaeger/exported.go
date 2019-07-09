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

		// Log tracer identifying values
		log.WithFields(log.Fields{
			"service.name": reflect.ValueOf(tracer).Elem().FieldByName("serviceName"),
			"service.tags": reflect.ValueOf(tracer).Elem().FieldByName("tags"),
		}).Infof("jaeger: tracer agent ready")

		// Set Open tracing global tracer
		opentracing.SetGlobalTracer(tracer)
	})
}
