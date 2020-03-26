package contractregistry

import (
	"context"
	"sync"

	"github.com/go-pg/pg/v9"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	grpcserver "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/server/grpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/server/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/server/rest"
	servicelayer "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/service"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"
)

var (
	app       = common.NewApp()
	startOnce = &sync.Once{}
)

// StartService Starts gRPC and HTTP servers
func StartService(ctx context.Context) {
	startOnce.Do(func() {
		log.Info("Starting service")

		cancelCtx, cancel := context.WithCancel(ctx)

		// Initialize dependencies
		multitenancy.Init(ctx)
		db := createDBConnection()

		// Start services
		contractRegistryController := initializeController(db)

		go metrics.StartServer(ctx, cancel, app.IsAlive, app.IsReady)
		go servicelayer.InitGRPC(cancelCtx, cancel, contractRegistryController)
		servicelayer.InitHTTP(cancelCtx)

		// Initialize contracts from command line
		initializeABIs(ctx, contractRegistryController)

		// Indicate that application is ready
		app.SetReady(true)

		servicelayer.ListenAndServe(cancel)
	})
}

// StopService gracefully stops the application
func StopService(ctx context.Context) {
	log.Warn("app: stopping...")
	app.SetReady(false)
	common.InParallel(
		func() { grpcserver.StopServer(ctx) },
		func() { metrics.StopServer(ctx) },
		func() { rest.StopServer(ctx) },
	)
	log.Info("app: gracefully stopped application")
}

// Flags set up necessary flags for the contract registry
func Flags(runCmd *cobra.Command) {
	log.Info("Setting flags")

	// Hostname & port for servers
	grpcserver.Flags(runCmd.Flags())
	rest.Flags(runCmd.Flags())

	// ContractRegistry flags
	bindTypeFlag(runCmd.Flags())
	bindABIFlag(runCmd.Flags())

	// Postgres flags
	postgres.PGFlags(runCmd.Flags())
}

func createDBConnection() *pg.DB {
	storeType := viper.GetString(typeViperKey)

	switch storeType {
	case postgresOpt:
		opts := postgres.NewOptions()

		log.WithFields(log.Fields{
			"db.address":  opts.Addr,
			"db.database": opts.Database,
			"db.user":     opts.User,
		}).Infof("Contract registry db connected")

		return postgres.New(opts)
	default:
		log.WithFields(log.Fields{
			"type": viper.GetString(typeViperKey),
		}).Fatalf("%s: unknown store type", storeType)

		return nil
	}
}
