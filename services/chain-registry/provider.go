package chainregistry

import (
	"context"
	"fmt"
	"math"
	"time"

	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider/aggregator"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider/poll"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/chain-registry/use-cases/chains"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

const (
	InternalProviderName = "internal"
	ChainsProxyProvider  = "chains-proxy"
)

func NewProvider(
	getChains usecases.GetChains,
	refresh time.Duration,
	proxyCacheTTL *time.Duration,
) provider.Provider {
	prvdr := aggregator.New()
	prvdr.AddProvider(NewInternalProvider())
	prvdr.AddProvider(NewChainsProxyProvider(getChains, refresh, proxyCacheTTL))
	return prvdr
}

func NewInternalProvider() provider.Provider {
	return static.New(dynamic.NewMessage(InternalProviderName, NewInternalConfig()))
}

func NewChainsProxyProvider(getChains usecases.GetChains, refresh time.Duration, proxyCacheTTL *time.Duration) provider.Provider {
	poller := func(ctx context.Context) (provider.Message, error) {
		chains, err := getChains.Execute(ctx, []string{multitenancy.Wildcard}, nil)
		if err != nil {
			return nil, err
		}

		return dynamic.NewMessage(ChainsProxyProvider, NewProxyConfig(chains, proxyCacheTTL)), nil
	}
	return poll.New(poller, refresh)
}

func NewInternalConfig() *dynamic.Configuration {
	dynamicCfg := dynamic.NewConfig()

	// Router to Chains API
	dynamicCfg.HTTP.Routers["chains"] = &dynamic.Router{
		Router: &traefikdynamic.Router{
			EntryPoints: []string{http.DefaultHTTPAppEntryPoint},
			Service:     "chains",
			Priority:    math.MaxInt32,
			Rule:        "PathPrefix(`/chains`) || PathPrefix(`/faucets`)",
			Middlewares: []string{"base@logger-base", "auth@multitenancy"},
		},
	}

	// Chains API
	dynamicCfg.HTTP.Services["chains"] = &dynamic.Service{
		Chains: &dynamic.Chains{},
	}

	// Log Middleware for Chains
	dynamicCfg.HTTP.Middlewares["chain-proxy-accesslog"] = &dynamic.Middleware{
		AccessLog: &dynamic.AccessLog{
			Filters: &dynamic.AccessLogFilters{
				StatusCodes: []string{"100-199", "400-428", "430-599"},
			},
		},
	}

	// Middleware used by Chain-Proxy
	dynamicCfg.HTTP.Middlewares["strip-path"] = &dynamic.Middleware{
		Middleware: &traefikdynamic.Middleware{
			StripPrefixRegex: &traefikdynamic.StripPrefixRegex{
				Regex: []string{`/(?:tessera/)?(?:[a-zA-Z\d-]*)/?`},
			},
		},
	}

	// Rate Limit middleware for Chain proxy
	dynamicCfg.HTTP.Middlewares["ratelimit"] = &dynamic.Middleware{
		RateLimit: &dynamic.RateLimit{
			MaxDelay:     time.Second,
			DefaultDelay: 30 * time.Second,
			Cooldown:     30 * time.Second,
		},
	}

	return dynamicCfg
}

func NewProxyConfig(chains []*models.Chain, proxyCacheTTL *time.Duration) *dynamic.Configuration {
	cfg := dynamic.NewConfig()

	for _, chain := range chains {
		multitenancyMid := fmt.Sprintf("auth-%v@multitenancy", chain.TenantID)
		middlewares := []string{
			"chain-proxy-accesslog@internal",
			"auth@multitenancy",
			multitenancyMid,
			"strip-path@internal",
		}

		cfg.HTTP.Middlewares[multitenancyMid] = &dynamic.Middleware{
			MultiTenancy: &dynamic.MultiTenancy{
				Tenant: chain.TenantID,
			},
		}

		if proxyCacheTTL != nil {
			httpCacheMid := fmt.Sprintf("%s@http-cache", chain.UUID)
			middlewares = append(middlewares, httpCacheMid)
			cfg.HTTP.Middlewares[httpCacheMid] = &dynamic.Middleware{
				HTTPCache: &dynamic.HTTPCache{
					TTL:       *proxyCacheTTL,
					KeySuffix: httpCacheGenerateChainKey(chain),
				},
			}
		}

		middlewares = append(middlewares, "ratelimit@internal")

		appendChainServices(cfg, chain, middlewares)
		appendTesseraPrivateTxServices(cfg, chain, middlewares)
	}

	return cfg
}

func appendChainServices(cfg *dynamic.Configuration, chain *models.Chain, middlewares []string) {
	chainService := fmt.Sprintf("chain-%v", chain.UUID)

	servers := make([]*dynamic.Server, 0)
	for _, chainURL := range chain.URLs {
		servers = append(servers, &dynamic.Server{
			URL: chainURL,
		})
	}

	cfg.HTTP.Routers[chainService] = &dynamic.Router{
		Router: &traefikdynamic.Router{
			EntryPoints: []string{http.DefaultHTTPAppEntryPoint},
			Priority:    math.MaxInt32,
			Service:     chainService,
			Rule:        fmt.Sprintf("Path(`/%s`)", chain.UUID),
			Middlewares: middlewares,
		},
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

func appendTesseraPrivateTxServices(cfg *dynamic.Configuration, chain *models.Chain, middlewares []string) {
	servers := make([]*dynamic.Server, 0)
	for _, privTxManager := range chain.PrivateTxManagers {
		if privTxManager.Type == utils.TesseraChainType {
			servers = append(servers, &dynamic.Server{
				URL: privTxManager.URL,
			})
		}
	}

	// Not servers identified
	if len(servers) == 0 {
		return
	}

	chainService := fmt.Sprintf("tessera-chain-%v", chain.UUID)
	cfg.HTTP.Routers[chainService] = &dynamic.Router{
		Router: &traefikdynamic.Router{
			EntryPoints: []string{http.DefaultHTTPAppEntryPoint},
			Priority:    math.MaxInt32,
			Service:     chainService,
			Rule:        fmt.Sprintf("PathPrefix(`/tessera/%s`)", chain.UUID),
			Middlewares: middlewares,
		},
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
