package swagger

import (
	"math"

	"github.com/consensys/orchestrate/pkg/toolkit/app/http"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/config/dynamic"
	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
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
