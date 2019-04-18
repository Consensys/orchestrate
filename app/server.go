package app

import (
	"fmt"
	"net/http"

	"github.com/julien-marchand/healthcheck"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// initServer creates server and attach it to application
func initServer(app *App) {
	// Create server
	server := &http.Server{
		Addr:    viper.GetString("http.hostname"),
		Handler: prepareHTTPRouter(app),
	}

	// Attach server on application
	app.server = server

	// Start Server
	go func() {
		log.Infof("server listening on %v", server.Addr)
		err := server.ListenAndServe()
		if err != nil {
			log.Errorf("server error: %v", err)
			// We encounter an issue with the server so we stop the application
			app.Close()
		}
	}()

	// Wait for app to be done and then close all connection
	go func() {
		<-app.Done()
		server.Close()
	}()
}

func prepareHTTPRouter(app *App) *http.ServeMux {
	// Create a metrics-exposing Handler for the Prometheus registry
	// The healthcheck related metrics will be prefixed with the provided namespace
	health := healthcheck.NewMetricsHandler(prometheus.DefaultRegisterer, "health")

	// Add a liveness check that always succeeds
	health.AddLivenessCheck("liveness-check", func() error {
		return nil
	})

	// Add a simple readiness check that always fails.
	health.AddReadinessCheck("readiness-check", func() error {
		if !app.Ready() {
			return fmt.Errorf("not ready")
		}
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
