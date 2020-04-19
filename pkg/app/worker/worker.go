package worker

import (
	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	prom "github.com/prometheus/client_golang/prometheus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	dynhandler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/healthcheck"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/prometheus"
	dynmid "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/dynamic"
	metricsmid "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/router"
	dynrouter "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/router/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics"
	metricsmulti "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/multi"
)

const (
	InternalProvider = "internal"
)

func New(cfg *app.Config, registry prom.Registerer) (*app.App, error) {
	// Create metrics registry
	reg := metricsmulti.New(cfg.Metrics)
	err := registry.Register(reg.Prometheus())
	if err != nil {
		return nil, err
	}

	// Create app
	return app.New(
		cfg,
		NewProvider(cfg.HTTP),
		NewHTTPBuilder(cfg.HTTP, reg.HTTP()),
		nil,
		reg,
	)
}

func NewHTTPBuilder(staticCfg *traefikstatic.Configuration, reg metrics.HTTP) router.Builder {
	builder := dynrouter.NewBuilder(staticCfg, nil)
	builder.Handler = dynhandler.NewBuilder()
	builder.Middleware = dynmid.NewBuilder()
	builder.Metrics = metricsmid.NewBuilder(reg)
	return builder
}

func NewProvider(staticCfg *traefikstatic.Configuration) provider.Provider {
	return static.New(dynamic.NewMessage(InternalProvider, NewInternalConfig()))
}

func NewInternalConfig() *dynamic.Configuration {
	cfg := dynamic.NewConfig()
	healthcheck.AddDynamicConfig(cfg)
	prometheus.AddDynamicConfig(cfg)
	return cfg
}
