package proxy

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/consensys/orchestrate/pkg/toolkit/app/http"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/config/dynamic"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/configwatcher/provider"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/configwatcher/provider/poll"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/utils"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
)

const ChainsProxyProvider = "chains-proxy"

func NewChainsProxyProvider(searchChains usecases.SearchChainsUseCase, refresh time.Duration, proxyCacheTTL *time.Duration) provider.Provider {
	poller := func(ctx context.Context) (provider.Message, error) {
		chains, err := searchChains.Execute(ctx, &entities.ChainFilters{}, []string{multitenancy.Wildcard})
		if err != nil {
			return nil, err
		}

		return dynamic.NewMessage(ChainsProxyProvider, NewProxyConfig(chains, proxyCacheTTL)), nil
	}
	return poll.New(poller, refresh)
}

func NewProxyConfig(chains []*entities.Chain, proxyCacheTTL *time.Duration) *dynamic.Configuration {
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

		if chain.PrivateTxManager != nil {
			appendTesseraPrivateTxServices(cfg, chain, middlewares)
		}
	}

	return cfg
}

func NewInternalConfig(dynamicCfg *dynamic.Configuration) *dynamic.Configuration {
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
				Regex: []string{`/proxy/chains/(?:tessera/)?(?:[a-zA-Z\d-]*)/?`},
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

func appendChainServices(cfg *dynamic.Configuration, chain *entities.Chain, middlewares []string) {
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
			Rule:        fmt.Sprintf("Path(`/proxy/chains/%s`)", chain.UUID),
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

func appendTesseraPrivateTxServices(cfg *dynamic.Configuration, chain *entities.Chain, middlewares []string) {
	servers := make([]*dynamic.Server, 0)
	if chain.PrivateTxManager.Type == entities.TesseraChainType {
		servers = append(servers, &dynamic.Server{
			URL: chain.PrivateTxManager.URL,
		})
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
			Rule:        fmt.Sprintf("PathPrefix(`/proxy/chains/tessera/%s`)", chain.UUID),
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
