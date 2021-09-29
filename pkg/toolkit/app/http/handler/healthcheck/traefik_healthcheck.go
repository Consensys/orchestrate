package healthcheck

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/consensys/orchestrate/pkg/toolkit/app/http/config/dynamic"
	"github.com/gorilla/mux"
)

type TraefikBuilder struct{}

func NewTraefikBuilder() *TraefikBuilder {
	return &TraefikBuilder{}
}

func (b *TraefikBuilder) Build(ctx context.Context, _ string, configuration interface{}, respModifier func(*http.Response) error) (http.Handler, error) {
	cfg, ok := configuration.(*dynamic.HealthCheck)
	if !ok {
		return nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	router := mux.NewRouter()
	h := NewTraefikHealthCheck()
	h.WithContext(ctx)
	h.Append(router)

	return router, nil
}

// TraefikHealthCheck expose ping routes.
type TraefikHealthCheck struct {
	terminating int32
}

func NewTraefikHealthCheck() *TraefikHealthCheck {
	return &TraefikHealthCheck{}
}

// WithContext causes the ping endpoint to serve non 200 responses.
func (h *TraefikHealthCheck) WithContext(ctx context.Context) {
	go func() {
		<-ctx.Done()
		atomic.StoreInt32(&h.terminating, 1)
	}()
}

// Append adds ping routes on a router.
func (h *TraefikHealthCheck) Append(router *mux.Router) {
	router.Methods(http.MethodGet, http.MethodHead).Path("/live").HandlerFunc(h.isAlive)
	router.Methods(http.MethodGet, http.MethodHead).Path("/ready").HandlerFunc(h.isReady)
}

func (h *TraefikHealthCheck) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	statusCode := http.StatusOK
	if atomic.LoadInt32(&h.terminating) > 0 {
		statusCode = http.StatusServiceUnavailable
	}
	response.WriteHeader(statusCode)
	_, _ = fmt.Fprint(response, http.StatusText(statusCode))
}

// IsAlive indicates if application is alive
func (h *TraefikHealthCheck) isAlive(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(response, http.StatusText(http.StatusOK))
}

// Ready indicates if application is ready
func (h *TraefikHealthCheck) isReady(response http.ResponseWriter, request *http.Request) {
	statusCode := http.StatusOK
	if atomic.LoadInt32(&h.terminating) > 0 {
		statusCode = http.StatusServiceUnavailable
	}
	response.WriteHeader(statusCode)
	_, _ = fmt.Fprint(response, http.StatusText(statusCode))
}
