package transactionscheduler

import (
	"context"
	"reflect"

	"github.com/Shopify/sarama"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/controllers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/multi"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases"
)

func New(
	cfg *Config,
	pgmngr postgres.Manager,
	jwt, key auth.Checker,
	chainRegistryClient client.ChainRegistryClient,
	syncProducer sarama.SyncProducer,
	txCrafterTopic string,
) (*app.App, error) {
	// Create Data agents
	db, err := multi.Build(context.Background(), cfg.Store, pgmngr)
	if err != nil {
		return nil, err
	}

	ucs := usecases.NewUseCases(db, chainRegistryClient, syncProducer, txCrafterTopic)
	// Option for transaction handler
	txSchedulerHandlerOpt := app.HandlerOpt(
		reflect.TypeOf(&dynamic.Transactions{}),
		controllers.NewBuilder(ucs),
	)

	// Create app
	return app.New(
		cfg.App,
		app.MultiTenancyOpt("auth", jwt, key, cfg.Multitenancy),
		app.MetricsOpt(),
		app.LoggerMiddlewareOpt("base"),
		app.SwaggerOpt("./public/swagger-specs/types/transaction-scheduler/swagger.json", "base@logger-base"),
		txSchedulerHandlerOpt,
		app.ProviderOpt(NewProvider()),
	)
}
