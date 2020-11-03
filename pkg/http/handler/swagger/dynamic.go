package swagger

import (
	"math"

	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
)

func AddDynamicConfig(cfg *dynamic.Configuration, middlewares []string, specsFile string) {
	// Router to swagger
	cfg.HTTP.Routers["swagger"] = &dynamic.Router{
		Router: &traefikdynamic.Router{
			EntryPoints: []string{http.DefaultHTTPAppEntryPoint},
			Service:     "swagger",
			Priority:    math.MaxInt32,
			Rule:        "PathPrefix(`/swagger`)",
			Middlewares: middlewares,
		},
	}

	// Swagger
	cfg.HTTP.Services["swagger"] = &dynamic.Service{
		Swagger: &dynamic.Swagger{
			SpecsFile: specsFile,
		},
	}
}
