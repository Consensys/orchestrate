package jaeger

import (
	"context"
	"reflect"
	"sync"

	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/rpcmetrics"
	jaegermetrics "github.com/uber/jaeger-lib/metrics"
	prometheus "github.com/uber/jaeger-lib/metrics/prometheus"
)

var (
	cfg      *jaegercfg.Configuration
	initOnce = &sync.Once{}
)

// Init initialize tracer
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if cfg == nil {
			cfg = NewConfig()
		}

		metrics := prometheus.New()
		tracer, _, err := cfg.NewTracer(
			jaegercfg.Logger(logger{entry: log.StandardLogger().WithFields(log.Fields{"system": "opentracing.jaeger"})}),
			jaegercfg.Observer(rpcmetrics.NewObserver(metrics.Namespace(jaegermetrics.NSOptions{Name: cfg.ServiceName}), rpcmetrics.DefaultNameNormalizer)),
		)
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
