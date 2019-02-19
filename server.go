package main

import (
	"context"
	"net/http"

	"github.com/heptiolabs/healthcheck"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func prepareHTTPRouter(ctx context.Context) *http.ServeMux {
	// Create a metrics-exposing Handler for the Prometheus registry
	// The healthcheck related metrics will be prefixed with the provided namespace
	health := healthcheck.NewMetricsHandler(prometheus.DefaultRegisterer, "health")

	// Add a liveness check that always succeeds
	health.AddLivenessCheck("liveness-check", func() error {
		return nil
	})

	// Add a simple readiness check that always fails.
	health.AddReadinessCheck("readiness-check", func() error {
		// TODO: return error if running and can access external components.
		return nil
	})

	router := http.NewServeMux()

	// Expose prometheus metrics on /metrics
	router.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))

	// Expose a liveness check on /live
	router.HandleFunc("/live", health.LiveEndpoint)

	// Expose a readiness check on /ready
	router.HandleFunc("/ready", health.ReadyEndpoint)

	// Return HTTP Server instance
	return router
}
