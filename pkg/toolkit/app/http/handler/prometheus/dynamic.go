package prometheus

import (
	"math"

	"github.com/consensys/orchestrate/pkg/toolkit/app/http"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/config/dynamic"
	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
)

func AddDynamicConfig(cfg *dynamic.Configuration) {
	// Router to Healthchecks
	cfg.HTTP.Routers["prometheus"] = &dynamic.Router{
		Router: &traefikdynamic.Router{
			EntryPoints: []string{http.DefaultMetricsEntryPoint},
			Service:     "prometheus",
			Priority:    math.MaxInt32,
			Rule:        "PathPrefix(`/metrics`)",
		},
	}

	// Healthcheck
	cfg.HTTP.Services["prometheus"] = &dynamic.Service{
		Prometheus: &dynamic.Prometheus{},
	}
}
