package chainregistry

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	authjwt "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/jwt"
	authkey "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/key"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient/rpc"
	orchlog "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
)

var (
	cfg      *app.Config
	appli    *app.App
	initOnce = &sync.Once{}
)

func initDependencies(ctx context.Context) {
	authjwt.Init(ctx)
	authkey.Init(ctx)
	rpc.Init(ctx)
}

func Init(ctx context.Context) {
	initOnce.Do(func() {
		if appli != nil {
			return
		}

		if cfg == nil {
			cfg = app.NewConfig(viper.GetViper())
		}
		orchlog.ConfigureLogger(cfg.HTTP)

		jsonConf, err := json.Marshal(cfg.HTTP)
		if err != nil {
			log.WithoutContext().WithError(err).Fatalf("could not marshal HTTP configuration: %#v", cfg.HTTP)
		} else {
			log.WithoutContext().Infof("HTTP configuration loaded %s", string(jsonConf))
		}

		// Initialize dependencied
		initDependencies(ctx)

		s, err := NewStore(postgres.GetManager(), store.NewConfig(viper.GetViper()))
		if err != nil {
			log.WithoutContext().WithError(err).Fatalf("could not create chain store")
		}

		// Init Chains
		store.ImportChains(ctx, s, viper.GetStringSlice(store.InitViperKey))

		appli, err = New(
			cfg,
			authjwt.GlobalChecker(),
			authkey.GlobalChecker(),
			viper.GetBool(multitenancy.EnabledViperKey),
			s,
			rpc.GlobalClient(),
		)
		if err != nil {
			log.FromContext(ctx).WithError(err).Fatalf("Could not create application")
		}
	})
}

func Start(ctx context.Context) error {
	Init(ctx)
	return appli.Start(ctx)
}

func Stop(ctx context.Context) error {
	return appli.Stop(ctx)
}
