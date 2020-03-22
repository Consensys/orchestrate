package healthcheck

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/julien-marchand/healthcheck"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
)

type Builder struct {
	health healthcheck.Handler
}

func NewBuilder(health healthcheck.Handler) *Builder {
	return &Builder{
		health: health,
	}
}

func (b *Builder) Health() healthcheck.Handler {
	return b.health
}

func (b *Builder) Build(ctx context.Context, _ string, configuration interface{}, respModifier func(*http.Response) error) (http.Handler, error) {
	cfg, ok := configuration.(*dynamic.HealthCheck)
	if !ok {
		return nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

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
	router.HandleFunc("/live", h.health.LiveEndpoint)
	router.HandleFunc("/ready", h.health.ReadyEndpoint)
}
