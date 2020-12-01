package keymanager

import (
	"math"

	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/configwatcher/provider"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/configwatcher/provider/aggregator"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/configwatcher/provider/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/config/dynamic"
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

	// Router to Key management API
	cfg.HTTP.Routers["ethereum"] = &dynamic.Router{
		Router: &traefikdynamic.Router{
			EntryPoints: []string{http.DefaultHTTPAppEntryPoint},
			Service:     "ethereum",
			Priority:    math.MaxInt32,
			Rule:        "PathPrefix(`/ethereum`)",
			Middlewares: []string{"base@logger-base"},
		},
	}

	// Ethereum accounts API
	cfg.HTTP.Services["ethereum"] = &dynamic.Service{
		Signer: &dynamic.Signer{},
	}

	return cfg
}
