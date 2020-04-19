package configwatcher

import (
	"fmt"
	"math"
	"time"

	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/dashboard"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/healthcheck"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/prometheus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/swagger"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/accesslog"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

type Config struct {
	static  *traefikstatic.Configuration
	watcher *configwatcher.Config
	dynamic *dynamic.Configuration
}

func NewInternalConfig(staticCfg *traefikstatic.Configuration, watcherCfg *configwatcher.Config) *Config {
	dynamicCfg := dynamic.NewConfig()

	dashboard.AddDynamicConfig(dynamicCfg, []string{"base-accesslog"})
	swagger.AddDynamicConfig(dynamicCfg,
		[]string{"base-accesslog"},
		"./public/swagger-specs/types/chain-registry/swagger.json",
	)
	healthcheck.AddDynamicConfig(dynamicCfg)
	prometheus.AddDynamicConfig(dynamicCfg)
	accesslog.AddDynamicConfig(dynamicCfg, "base-accesslog", staticCfg)

	// Authentication middleware
	dynamicCfg.HTTP.Middlewares["auth"] = &dynamic.Middleware{
		Auth: &dynamic.Auth{},
	}

	// Router to Chains API
	dynamicCfg.HTTP.Routers["chains"] = &dynamic.Router{
		Router: &traefikdynamic.Router{
			EntryPoints: []string{http.DefaultHTTPEntryPoint},
			Service:     "chains",
			Priority:    math.MaxInt32,
			Rule:        "PathPrefix(`/chains`) || PathPrefix(`/faucets`)",
			Middlewares: []string{"base-accesslog", "auth"},
		},
	}

	// Chains API
	dynamicCfg.HTTP.Services["chains"] = &dynamic.Service{
		Chains: &dynamic.Chains{},
	}

	accesslog.AddDynamicConfig(dynamicCfg, "chain-proxy-accesslog", staticCfg)
	dynamicCfg.HTTP.Middlewares["chain-proxy-accesslog"].AccessLog.Filters = &dynamic.AccessLogFilters{
		StatusCodes: []string{"100-199", "400-428", "430-599"},
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

	return &Config{
		static:  staticCfg,
		watcher: watcherCfg,
		dynamic: dynamicCfg,
	}
}

func (c *Config) DynamicCfg() *dynamic.Configuration {
	return c.dynamic
}

func newProxyConfig(chains []*models.Chain) *dynamic.Configuration {
	cfg := dynamic.NewConfig()

	for _, chain := range chains {
		multitenancyMid := fmt.Sprintf("multitenancy-%v", chain.TenantID)
		middlewares := []string{
			"chain-proxy-accesslog@internal",
			"auth@internal",
			multitenancyMid,
			"strip-path@internal",
			"ratelimit@internal",
		}

		cfg.HTTP.Middlewares[multitenancyMid] = &dynamic.Middleware{
			MultiTenancy: &dynamic.MultiTenancy{
				Tenant: chain.TenantID,
			},
		}

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
			EntryPoints: []string{http.DefaultHTTPEntryPoint},
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
			EntryPoints: []string{http.DefaultHTTPEntryPoint},
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
