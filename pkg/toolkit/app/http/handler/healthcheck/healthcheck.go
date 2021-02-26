package healthcheck

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/heptiolabs/healthcheck"
)

type Checker struct {
	Name  string
	Check healthcheck.Check
}

func NewChecker(name string, check healthcheck.Check) *Checker {
	return &Checker{
		Name:  name,
		Check: check,
	}
}

type Builder struct {
	health healthcheck.Handler
}

func NewBuilder(handler healthcheck.Handler) *Builder {
	return &Builder{
		health: handler,
	}
}

func (b *Builder) Health() healthcheck.Handler {
	return b.health
}

func (b *Builder) Build(_ context.Context, _ string, _ interface{}, _ func(*http.Response) error) (http.Handler, error) {
	router := mux.NewRouter()
	New(b.health).Append(router)

	return router, nil
}

type HealthCheck struct {
	health healthcheck.Handler
}

func New(health healthcheck.Handler) *HealthCheck {
	return &HealthCheck{
		health: health,
	}
}

// Append add dashboard routes on a router
func (h *HealthCheck) Append(router *mux.Router) {
	router.HandleFunc("/live", h.isAlive)
	router.HandleFunc("/ready", h.health.ReadyEndpoint)
}

// IsAlive indicates if application is alive
func (h *HealthCheck) isAlive(response http.ResponseWriter, _ *http.Request) {
	response.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(response, http.StatusText(http.StatusOK))
}
