package contractregistry

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	authjwt "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/jwt"
	authkey "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/key"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	orchlog "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/logger"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store"
)

var (
	cfg      *app.Config
	appli    *app.App
	initOnce = &sync.Once{}
)

func initDependencies(ctx context.Context) {
	authjwt.Init(ctx)
	authkey.Init(ctx)
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

		// Create GRPC service
		service, err := NewService(postgres.GetManager(), store.NewConfig(viper.GetViper()))
		if err != nil {
			log.WithoutContext().WithError(err).Fatalf("could not create contracts GRPC service")
		}

		// Create App
		appli, err = New(
			cfg,
			authjwt.GlobalChecker(),
			authkey.GlobalChecker(),
			viper.GetBool(multitenancy.EnabledViperKey),
			service,
			logrus.StandardLogger(),
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

// // Flags set up necessary flags for the contract registry
// func Flags(runCmd *cobra.Command) {
// 	log.Info("Setting flags")

// 	// Hostname & port for servers
// 	grpcserver.Flags(runCmd.Flags())
// 	rest.Flags(runCmd.Flags())

// 	// ContractRegistry flags
// 	bindTypeFlag(runCmd.Flags())
// 	bindABIFlag(runCmd.Flags())

// 	// Postgres flags
// 	postgres.PGFlags(runCmd.Flags())
// }
