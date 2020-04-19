package configwatcher

import (
	"math"

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
)

type Config struct {
	static  *traefikstatic.Configuration
	watcher *configwatcher.Config
	dynamic *dynamic.Configuration
}

func NewConfig(staticCfg *traefikstatic.Configuration, watcherCfg *configwatcher.Config) *Config {
	dynamicCfg := dynamic.NewConfig()

	dashboard.AddDynamicConfig(dynamicCfg, []string{"base-accesslog"})
	swagger.AddDynamicConfig(dynamicCfg,
		[]string{"base-accesslog"},
		"./public/swagger-specs/types/transaction-scheduler/swagger.json",
	)
	healthcheck.AddDynamicConfig(dynamicCfg)
	prometheus.AddDynamicConfig(dynamicCfg)
	accesslog.AddDynamicConfig(dynamicCfg, "base-accesslog", staticCfg)

	// Authentication middleware
	dynamicCfg.HTTP.Middlewares["auth"] = &dynamic.Middleware{
		Auth: &dynamic.Auth{},
	}

	// Router to Transactions API
	dynamicCfg.HTTP.Routers["transactions"] = &dynamic.Router{
		Router: &traefikdynamic.Router{
			EntryPoints: []string{http.DefaultHTTPEntryPoint},
			Service:     "transactions",
			Priority:    math.MaxInt32,
			Rule:        "PathPrefix(`/transactions`)",
			Middlewares: []string{"base-accesslog", "auth"},
		},
	}

	// Transaction scheduler API
	dynamicCfg.HTTP.Services["transactions"] = &dynamic.Service{
		Transactions: &dynamic.Transactions{},
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
