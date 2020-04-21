package chainregistry

import (
	"context"

	"github.com/containous/traefik/v2/pkg/log"
	pkgapp "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	authjwt "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/jwt"
	authkey "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/key"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient/rpc"
	orchlog "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/log"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/chain-registry/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
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
	rpc.Init(ctx)

	storeBuilder := store.NewBuilder(postgres.GetManager())
	storeDataAgents, err := storeBuilder.Build(context.Background(), cfg.store)
	if err != nil {
		logger.WithError(err).Fatalf("could not create data-agents")
		return nil, err
	}

	// Init Chains
	importChainUC := usecases.NewImportChain(storeDataAgents.Chain, rpc.GlobalClient())
	for _, jsonChain := range cfg.envChains {
		err = importChainUC.Execute(ctx, jsonChain)
		if err != nil {
			logger.WithError(err).Errorf("could not import chain")
		}
	}

	app, err := newApplication(
		cfg,
		storeDataAgents,
		rpc.GlobalClient(),
		authjwt.GlobalChecker(),
		authkey.GlobalChecker(),
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
