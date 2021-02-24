package dashboard

import (
	"math"

	"github.com/ConsenSys/orchestrate/pkg/http"
	"github.com/ConsenSys/orchestrate/pkg/http/config/dynamic"
	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
)

func AddDynamicConfig(cfg *dynamic.Configuration, middlewares []string) {
	cfg.HTTP.Routers["dashboard"] = &dynamic.Router{
		Router: &traefikdynamic.Router{
			EntryPoints: []string{http.DefaultHTTPAppEntryPoint},
			Service:     "dashboard",
			Priority:    math.MaxInt32,
			Rule:        "PathPrefix(`/api`) || PathPrefix(`/dashboard`)",
			Middlewares: append(middlewares, "strip-api"),
		},
	}

	cfg.HTTP.Middlewares["strip-api"] = &dynamic.Middleware{
		Middleware: &traefikdynamic.Middleware{
			StripPrefixRegex: &traefikdynamic.StripPrefixRegex{
				Regex: []string{"/api"},
			},
		},
	}

	// Dashboard API
	cfg.HTTP.Services["dashboard"] = &dynamic.Service{
		Dashboard: &dynamic.Dashboard{},
	}
}
