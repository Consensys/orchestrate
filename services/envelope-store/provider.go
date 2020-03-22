package envelopestore

import (
	"math"

	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider/aggregator"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/dashboard"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/healthcheck"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/swagger"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/accesslog"
)

const (
	InternalProvider = "internal"
)

func NewProvider(
	staticCfg *traefikstatic.Configuration,
) provider.Provider {
	prvdr := aggregator.New()
	prvdr.AddProvider(NewInternalProvider(staticCfg))
	return prvdr
}

func NewInternalProvider(staticCfg *traefikstatic.Configuration) provider.Provider {
	return static.New(dynamic.NewMessage(InternalProvider, NewInternalConfig(staticCfg)))
}

func NewInternalConfig(staticCfg *traefikstatic.Configuration) *dynamic.Configuration {
	cfg := dynamic.NewConfig()

	dashboard.AddDynamicConfig(cfg, []string{"base-accesslog"})
	swagger.AddDynamicConfig(cfg,
		[]string{"base-accesslog"},
		"./public/swagger-specs/types/envelope-store/store.swagger.json",
	)
	healthcheck.AddDynamicConfig(cfg)
	accesslog.AddDynamicConfig(cfg, "base-accesslog", staticCfg)
	// Authentication middleware
	cfg.HTTP.Middlewares["auth"] = &dynamic.Middleware{
		Auth: &dynamic.Auth{},
	}

	// Router to Chains API
	cfg.HTTP.Routers["envelopes"] = &dynamic.Router{
		Router: &traefikdynamic.Router{
			EntryPoints: []string{http.DefaultHTTPEntryPoint},
			Service:     "envelopes",
			Priority:    math.MaxInt32,
			Rule:        "PathPrefix(`/envelopes`)",
			Middlewares: []string{"base-accesslog", "auth"},
		},
	}

	// Envelope API
	cfg.HTTP.Services["envelopes"] = &dynamic.Service{
		Envelopes: &dynamic.Envelopes{},
	}

	return cfg
}
