package prometheus

import (
	"math"

	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/config/dynamic"
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
