package worker

import (
	"context"

	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	dynhandler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/healthcheck"
	dynmid "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/router"
	dynrouter "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/router/dynamic"
)

const (
	InternalProvider = "internal"
)

func New(cfg *app.Config) *app.App {
	// Create HTTP Router builder
	httpBuilder := NewHTTPBuilder(cfg.HTTP)

	// Create HTTP EntryPoints
	httpEps := http.NewEntryPoints(
		cfg.HTTP.EntryPoints,
		httpBuilder,
	)

	// Create Configuration Watcher
	// Create configuration listener switching HTTP Endpoints configuration
	listeners := []func(context.Context, interface{}) error{
		httpEps.Switch,
	}

	watcher := configwatcher.New(
		cfg.Watcher,
		NewProvider(cfg.HTTP),
		dynamic.Merge,
		listeners,
	)

	// Create app
	return app.New(watcher, httpEps, nil)
}

func NewHTTPBuilder(staticCfg *traefikstatic.Configuration) router.Builder {
	builder := dynrouter.NewBuilder(staticCfg, nil)
	builder.Handler = dynhandler.NewBuilder()
	builder.Middleware = dynmid.NewBuilder()
	return builder
}

func NewProvider(staticCfg *traefikstatic.Configuration) provider.Provider {
	return static.New(dynamic.NewMessage(InternalProvider, NewInternalConfig()))
}

func NewInternalConfig() *dynamic.Configuration {
	cfg := dynamic.NewConfig()
	healthcheck.AddDynamicConfig(cfg)
	return cfg
}
