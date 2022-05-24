package api

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/consensys/orchestrate/pkg/toolkit/app/http"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/config/dynamic"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/configwatcher/provider"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/configwatcher/provider/aggregator"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/configwatcher/provider/static"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
	"github.com/consensys/orchestrate/services/api/proxy"
	traefikdynamic "github.com/traefik/traefik/v2/pkg/config/dynamic"
)

const (
	InternalProvider = "internal"
)

func NewProvider(searchChains usecases.SearchChainsUseCase, refresh time.Duration, proxyCacheTTL *time.Duration, accessLog bool) provider.Provider {
	prvdr := aggregator.New()
	prvdr.AddProvider(NewInternalProvider())
	prvdr.AddProvider(proxy.NewChainsProxyProvider(searchChains, refresh, proxyCacheTTL, accessLog))
	return prvdr

}

func NewInternalProvider() provider.Provider {
	return static.New(dynamic.NewMessage(InternalProvider, newInternalConfig()))
}

func newInternalConfig() *dynamic.Configuration {
	cfg := dynamic.NewConfig()

	pathPrefix := []string{"/transactions", "/schedules", "/jobs", "/accounts", "/faucets", "/contracts", "/chains"}
	for idx, path := range pathPrefix {
		pathPrefix[idx] = fmt.Sprintf("PathPrefix(`%s`)", path)
	}

	// Router to API
	cfg.HTTP.Routers["api"] = &dynamic.Router{
		Router: &traefikdynamic.Router{
			EntryPoints: []string{http.DefaultHTTPAppEntryPoint},
			Service:     "api",
			Priority:    math.MaxInt32,
			Rule:        strings.Join(pathPrefix, " || "),
			Middlewares: []string{"base@logger-base", "auth@multitenancy"},
		},
	}

	// API
	cfg.HTTP.Services["api"] = &dynamic.Service{
		API: &dynamic.API{},
	}

	return proxy.NewInternalConfig(cfg)
}
