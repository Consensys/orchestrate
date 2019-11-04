package healthcheck

import (
	"net/http"

	"github.com/julien-marchand/healthcheck"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	server "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
)

// App interface
type App interface {
	IsReady() error
}

// HealthCheck register HTTP handlers for application healthcheck
// TODO: still look quite opinionated and not very configurable
func HealthCheck(app App) server.ServeMuxEnhancer {
	return func(mux *http.ServeMux) *http.ServeMux {
		// Create a metrics-exposing Handler for the Prometheus registry
		// The healthcheck related metrics will be prefixed with the provided namespace
		health := healthcheck.NewMetricsHandler(prometheus.DefaultRegisterer, "health")

		// Add a liveness check that always succeeds
		health.AddLivenessCheck("liveness-check", func() error { return nil })

		// Add a simple readiness check that always fails.
		health.AddReadinessCheck("readiness-check", app.IsReady)

		// Expose prometheus metrics on /metrics
		mux.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))

		// Expose a liveness check on /live
		mux.HandleFunc("/live", health.LiveEndpoint)

		// Expose a readiness check on /ready
		mux.HandleFunc("/ready", health.ReadyEndpoint)

		// Return mux
		return mux
	}
}
