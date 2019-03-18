package http

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/heptiolabs/healthcheck"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/app/grpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/app/infra"
)

// Variable use as HTTP server singleton for injection pattern
var s *server

func init() {
	s = new()
}

type server struct {
	http *http.Server

	initOnce, closeOnce *sync.Once
	ready               *atomic.Value
}

// new create a new server
func new() *server {
	return &server{
		initOnce:  &sync.Once{},
		closeOnce: &sync.Once{},
		ready:     &atomic.Value{},
	}
}

// Init inititalize grpc server
func Init() {
	s.initOnce.Do(func() {
		s.http = CreateServer()
		s.ready.Store(true)
		log.Infof("http: server ready")
	})
}

// Ready indicate if server is ready
func Ready() bool {
	return s.ready.Load().(bool)
}

// Server returns HTTP server
func Server() *http.Server {
	if !Ready() {
		panic("GRPC server is not ready. Please call Init() first")
	}
	return s.http
}

// Close http server
func Close(ctx context.Context) {
	log.Debugf("http: closing...")
	Server().Shutdown(ctx)
	log.Debugf("http: closed")
}

// RegisterHealthChecks register health check handlers
func RegisterHealthChecks(router *http.ServeMux) {
	// Create a metrics-exposing Handler for the Prometheus registry
	// The healthcheck related metrics will be prefixed with the provided namespace
	health := healthcheck.NewMetricsHandler(prometheus.DefaultRegisterer, "health")

	// Add a liveness check that always succeeds
	health.AddLivenessCheck("liveness-check", func() error {
		return nil
	})

	// Add a simple readiness check that always fails.
	health.AddReadinessCheck("readiness-check", func() error {
		if !infra.Ready() || !grpc.Ready() {
			return fmt.Errorf("App is not ready")
		}
		return nil
	})

	// Expose prometheus metrics on /metrics
	router.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))

	// Expose a liveness check on /live
	router.HandleFunc("/live", health.LiveEndpoint)

	// Expose a readiness check on /ready
	router.HandleFunc("/ready", health.ReadyEndpoint)
}

// CreateServer create http server for the application
func CreateServer() *http.Server {
	router := http.NewServeMux()

	// Register healthcheck handlers
	RegisterHealthChecks(router)

	server := &http.Server{
		Handler: router,
	}

	// Return HTTP Server instance
	return server
}
