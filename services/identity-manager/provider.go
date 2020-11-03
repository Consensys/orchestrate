package identitymanager

import (
	"math"

	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider/aggregator"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
)

const (
	InternalProvider = "internal"
)

func NewProvider() provider.Provider {
	prvdr := aggregator.New()
	prvdr.AddProvider(NewInternalProvider())
	return prvdr
}

func NewInternalProvider() provider.Provider {
	return static.New(dynamic.NewMessage(InternalProvider, NewInternalConfig()))
}

func NewInternalConfig() *dynamic.Configuration {
	cfg := dynamic.NewConfig()

	// Router to Account API
	cfg.HTTP.Routers["identities"] = &dynamic.Router{
		Router: &traefikdynamic.Router{
			EntryPoints: []string{http.DefaultHTTPAppEntryPoint},
			Service:     "identities",
			Priority:    math.MaxInt32,
			Rule:        "PathPrefix(`/accounts`)",
			Middlewares: []string{"base@logger-base", "auth@multitenancy"},
		},
	}

	// Account API
	cfg.HTTP.Services["identities"] = &dynamic.Service{
		Identity: &dynamic.Identity{},
	}

	return cfg
}
