package opentracing

import (
	"context"
	"reflect"
	"sync"

	extOpentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/tracing/opentracing/jaeger"
)

var (
	initOnce     = &sync.Once{}
	globalTracer *Tracer
)

// Init initialize tracer
func Init(ctx context.Context) {
	initOnce.Do(func() {
		tracer := NewTracer(jaeger.TracerFromViperConfig())

		// Log tracer identifying values
		log.WithFields(log.Fields{
			"service.name": reflect.ValueOf(tracer.Internal).Elem().FieldByName("serviceName"),
			"service.tags": reflect.ValueOf(tracer.Internal).Elem().FieldByName("tags"),
		}).Infof("jaeger: tracer agent ready")

		SetGlobalTracer(tracer)

		// Also sets the global opentracing.Tracer for other services who don't interact with txctx
		extOpentracing.SetGlobalTracer(tracer.Internal)
	})
}

// SetGlobalTracer sets the global tracer instance
func SetGlobalTracer(tracer *Tracer) {
	globalTracer = tracer
}

// GetGlobalTracer returns the global tracer instance
func GetGlobalTracer() *Tracer {
	return globalTracer
}
