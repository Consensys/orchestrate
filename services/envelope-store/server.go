package envelopestore

import (
	"context"

	"github.com/containous/traefik/v2/pkg/log"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	pkgapp "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	authjwt "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/jwt"
	authkey "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/key"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	orchlog "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/logger"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/service/controllers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store"
)

type EnvelopStoreServer struct {
	ctx    context.Context
	cfg    *Config
	app    *pkgapp.App
	logger log.Logger
}

func NewServer(ctx context.Context, cfg *Config) (*EnvelopStoreServer, error) {
	logger := log.FromContext(ctx)

	orchlog.ConfigureLogger(cfg.app.HTTP)
	jsonConf, err := json.Marshal(cfg.app.HTTP)
	if err != nil {
		logger.WithError(err).Fatalf("could not marshal HTTP configuration: %#v", cfg.app.HTTP)
	} else {
		logger.Infof("HTTP configuration loaded %s", string(jsonConf))
	}

	// Initialize dependencies
	authjwt.Init(ctx)
	authkey.Init(ctx)

	// Create GRPC service
	pgmngr := postgres.GetManager()
	storeBuilder := store.NewBuilder(pgmngr)
	storeDataAgents, err := storeBuilder.Build(ctx, cfg.store)
	if err != nil {
		logger.WithError(err).Fatalf("could not create data-agents")
		return nil, err
	}

	srv, err := controllers.NewGRPCService(storeDataAgents)
	if err != nil {
		logger.WithError(err).Fatalf("could not create envelope store service")
		return nil, err
	}

	app, err := newApplication(
		cfg,
		authjwt.GlobalChecker(),
		authkey.GlobalChecker(),
		srv,
		logrus.StandardLogger(),
		prom.DefaultRegisterer,
	)

	if err != nil {
		logger.WithError(err).Fatalf("Could not create application")
		return nil, err
	}

	return &EnvelopStoreServer{
		ctx:    ctx,
		cfg:    cfg,
		app:    app,
		logger: logger,
	}, nil
}

func (serv *EnvelopStoreServer) Start() error {
	return serv.app.Start(serv.ctx)
}

func (serv *EnvelopStoreServer) Stop() error {
	return serv.app.Stop(serv.ctx)
}
