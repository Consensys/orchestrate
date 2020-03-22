package chainregistry

import (
	"context"
	"fmt"
	"math"
	"time"

	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider/aggregator"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider/poll"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/dashboard"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/healthcheck"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/swagger"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/accesslog"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

const (
	InternalProvider    = "internal"
	ChainsProxyProvider = "chains-proxy"
)

func NewProvider(
	staticCfg *traefikstatic.Configuration,
	s store.ChainRegistryStore,
) provider.Provider {
	prvdr := aggregator.New()
	prvdr.AddProvider(NewInternalProvider(staticCfg))
	prvdr.AddProvider(NewChainsProxyProvider(s, time.Second))
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
		"./public/swagger-specs/types/chain-registry/swagger.json",
	)
	healthcheck.AddDynamicConfig(cfg)
	accesslog.AddDynamicConfig(cfg, "base-accesslog", staticCfg)

	// Authentication middleware
	cfg.HTTP.Middlewares["auth"] = &dynamic.Middleware{
		Auth: &dynamic.Auth{},
	}

	// Router to Chains API
	cfg.HTTP.Routers["chains"] = &dynamic.Router{
		Router: &traefikdynamic.Router{
			EntryPoints: []string{http.DefaultHTTPEntryPoint},
			Service:     "chains",
			Priority:    math.MaxInt32,
			Rule:        "PathPrefix(`/chains`) || PathPrefix(`/faucets`)",
			Middlewares: []string{"base-accesslog", "auth"},
		},
	}

	// Chains API
	cfg.HTTP.Services["chains"] = &dynamic.Service{
		Chains: &dynamic.Chains{},
	}

	accesslog.AddDynamicConfig(cfg, "chain-proxy-accesslog", staticCfg)
	cfg.HTTP.Middlewares["chain-proxy-accesslog"].AccessLog.Filters = &dynamic.AccessLogFilters{
		StatusCodes: []string{"100-199", "400-428", "430-599"},
	}

	// Middleware used by Chain-Proxy
	cfg.HTTP.Middlewares["strip-path"] = &dynamic.Middleware{
		Middleware: &traefikdynamic.Middleware{
			StripPrefixRegex: &traefikdynamic.StripPrefixRegex{
				Regex: []string{"/.+"},
			},
		},
	}

	// Rate Limit middleware for Chain proxy
	cfg.HTTP.Middlewares["ratelimit"] = &dynamic.Middleware{
		RateLimit: &dynamic.RateLimit{
			MaxDelay:     time.Second,
			DefaultDelay: 30 * time.Second,
			Cooldown:     30 * time.Second,
		},
	}

	return cfg
}

func NewChainsProxyProvider(s store.ChainRegistryStore, refresh time.Duration) provider.Provider {
	poller := func(ctx context.Context) (provider.Message, error) {
		chains, err := s.GetChains(ctx, nil)
		if err != nil {
			return nil, err
		}

		return dynamic.NewMessage(ChainsProxyProvider, NewChainsProxyConfig(chains)), nil
	}
	return poll.New(poller, refresh)
}

func NewChainsProxyConfig(chains []*types.Chain) *dynamic.Configuration {
	cfg := dynamic.NewConfig()

	for _, chain := range chains {
		chainService := fmt.Sprintf("chain-%v", chain.UUID)
		multitenancyMid := fmt.Sprintf("multitenancy-%v", chain.TenantID)

		cfg.HTTP.Routers[chainService] = &dynamic.Router{
			Router: &traefikdynamic.Router{
				EntryPoints: []string{http.DefaultHTTPEntryPoint},
				Priority:    math.MaxInt32,
				Service:     chainService,
				Rule:        fmt.Sprintf("Path(`/%s`)", chain.UUID),
				Middlewares: []string{
					"chain-proxy-accesslog@internal",
					"auth@internal",
					multitenancyMid,
					"strip-path@internal",
					"ratelimit@internal",
				},
			},
		}

		cfg.HTTP.Middlewares[multitenancyMid] = &dynamic.Middleware{
			MultiTenancy: &dynamic.MultiTenancy{
				Tenant: chain.TenantID,
			},
		}

		servers := make([]*dynamic.Server, 0)
		for _, chainURL := range chain.URLs {
			servers = append(servers, &dynamic.Server{
				URL: chainURL,
			})
		}

		cfg.HTTP.Services[chainService] = &dynamic.Service{
			ReverseProxy: &dynamic.ReverseProxy{
				PassHostHeader: utils.Bool(false),
				LoadBalancer: &dynamic.LoadBalancer{
					Servers: servers,
				},
			},
		}
	}

	return cfg
}
