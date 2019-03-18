package infra

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/rpcmetrics"
	prometheus "github.com/uber/jaeger-lib/metrics/prometheus"
)

// logger is an adpatator from logrus to jaegger.Logger
type logger struct {
	log *log.Logger
}

func (l logger) Error(msg string) {
	l.log.Error(msg)
}

func (l logger) Infof(msg string, args ...interface{}) {
	l.log.Infof(msg, args...)
}

// InitTracing initilialize tracer
func InitTracing(infra *Infra) {
	metrics := prometheus.New()

	// Create tracer
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "const",
			Param: viper.GetFloat64("jaegger.sampler"),
		},
		Reporter: &jaegercfg.ReporterConfig{
			LocalAgentHostPort: fmt.Sprintf("%v:%d", viper.GetString("jaegger.host"), viper.GetInt("jaegger.port")),
		},
	}
	tracer, closer, err := cfg.New(
		"context-store",
		jaegercfg.Logger(logger{log: log.StandardLogger()}),
		jaegercfg.Observer(rpcmetrics.NewObserver(metrics.Namespace("context-store", nil), rpcmetrics.DefaultNameNormalizer)),
	)

	if err != nil {
		log.WithError(err).Fatalf("infra-tracing: could not initialize tracer")
	}

	// Set tracer
	infra.tracer = tracer
	log.Infof("infra-tracing: tracer ready")

	// Wait for app to be done and then close
	go func() {
		<-infra.ctx.Done()
		closer.Close()
	}()
}
