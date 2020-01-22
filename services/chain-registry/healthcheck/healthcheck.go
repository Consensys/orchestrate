package healthcheck

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Handler expose ping routes.
type Handler struct {
	EntryPoint  string `description:"EntryPoint" export:"true" json:"entryPoint,omitempty" toml:"entryPoint,omitempty" yaml:"entryPoint,omitempty"`
	terminating bool
}

// SetDefaults sets the default values.
func (h *Handler) SetDefaults() {
	h.EntryPoint = "metrics"
}

// WithContext causes the ping endpoint to serve non 200 responses.
func (h *Handler) WithContext(ctx context.Context) {
	go func() {
		<-ctx.Done()
		h.terminating = true
	}()
}

// Append adds ping routes on a router.
func (h *Handler) Append(router *mux.Router) {
	router.Methods(http.MethodGet, http.MethodHead).Path("/live").HandlerFunc(h.isAlive)
	router.Methods(http.MethodGet, http.MethodHead).Path("/ready").HandlerFunc(h.isReady)
}

func (h *Handler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	statusCode := http.StatusOK
	if h.terminating {
		statusCode = http.StatusServiceUnavailable
	}
	response.WriteHeader(statusCode)
	_, _ = fmt.Fprint(response, http.StatusText(statusCode))
}

// IsAlive indicates if application is alive
func (h *Handler) isAlive(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(response, http.StatusText(http.StatusOK))
}

// Ready indicates if application is ready
func (h *Handler) isReady(response http.ResponseWriter, request *http.Request) {
	statusCode := http.StatusOK
	if h.terminating {
		statusCode = http.StatusServiceUnavailable
	}
	response.WriteHeader(statusCode)
	_, _ = fmt.Fprint(response, http.StatusText(statusCode))
}
