package envelopestore

import (
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth"
	metrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/multi"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/service/configwatcher"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/service/grpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/service/http"
)

func newApplication(
	cfg *Config,
	jwt, key auth.Checker,
	srv svc.EnvelopeStoreServer,
	logger *logrus.Logger,
	registry prom.Registerer,
) (*app.App, error) {
	// Create metrics registry and register it on prometheus
	reg := metrics.New(cfg.app.Metrics)
	err := registry.Register(reg.Prometheus())
	if err != nil {
		return nil, err
	}

	// Create GRPC builder
	checker := auth.CombineCheckers(key, jwt)
	grpcBuilder, err := grpc.NewServerBuilder(srv, checker, cfg.multitenancy, logger, reg.GRPCServer())
	if err != nil {
		return nil, err
	}
	cfg.app.GRPC.Static = grpc.NewStaticConfig()

	// Create HTTP Router builder
	routerBuilder, err := http.NewRouterBuilder(srv, cfg.app.HTTP, jwt, key, cfg.multitenancy, reg.HTTP())
	if err != nil {
		return nil, err
	}

	// Create app
	return app.New(
		cfg.app,
		configwatcher.NewProvider(configwatcher.NewConfig(cfg.app.HTTP, cfg.app.Watcher).DynamicCfg()),
		routerBuilder,
		grpcBuilder,
		reg,
	)
}
