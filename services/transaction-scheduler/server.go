package transactionscheduler

import (
	"context"
	"encoding/json"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/sirupsen/logrus"
	pkgapp "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	authjwt "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/jwt"
	authkey "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/key"
	orchlog "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/logger"
)

type TransactionManagerServer struct {
	ctx    context.Context
	cfg    *Config
	app    *pkgapp.App
	logger log.Logger
}

func NewServer(ctx context.Context, cfg *Config) (*TransactionManagerServer, error) {
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
	client.Init(ctx)

	app, err := newApplication(
		ctx,
		cfg,
		authjwt.GlobalChecker(),
		authkey.GlobalChecker(),
		logrus.StandardLogger(),
		client.GlobalClient(),
	)

	if err != nil {
		logger.WithError(err).Fatalf("Could not create application")
		return nil, err
	}

	return &TransactionManagerServer{
		ctx:    ctx,
		cfg:    cfg,
		app:    app,
		logger: logger,
	}, nil
}

func (serv *TransactionManagerServer) Start() error {
	return serv.app.Start(serv.ctx)
}

func (serv *TransactionManagerServer) Stop() error {
	return serv.app.Stop(serv.ctx)
}
