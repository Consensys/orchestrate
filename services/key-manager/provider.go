package keymanager

import (
	"math"

	"github.com/ConsenSys/orchestrate/pkg/configwatcher/provider"
	"github.com/ConsenSys/orchestrate/pkg/configwatcher/provider/aggregator"
	"github.com/ConsenSys/orchestrate/pkg/configwatcher/provider/static"
	"github.com/ConsenSys/orchestrate/pkg/http"
	"github.com/ConsenSys/orchestrate/pkg/http/config/dynamic"
	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
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
	cfg.HTTP.Routers["key-manager"] = &dynamic.Router{
		Router: &traefikdynamic.Router{
			EntryPoints: []string{http.DefaultHTTPAppEntryPoint},
			Service:     "key-manager",
			Priority:    math.MaxInt32,
			Rule:        "PathPrefix(`/ethereum`) || PathPrefix(`/zk-snarks`)",
			Middlewares: []string{"base@logger-base"},
		},
	}

	// Ethereum accounts API
	cfg.HTTP.Services["key-manager"] = &dynamic.Service{
		KeyManager: &dynamic.KeyManager{},
	}

	return cfg
}
