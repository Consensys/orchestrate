package transactionscheduler

import (
	"context"
	"reflect"

	"github.com/Shopify/sarama"
	"github.com/go-pg/pg/v9/orm"
	pkgsarama "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database"
	contractregistry2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/client"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/builder"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	identitymanager "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/controllers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/multi"
)

func NewTxScheduler(
	cfg *Config,
	pgmngr postgres.Manager,
	jwt, key auth.Checker,
	chainRegistryClient chainregistry.ChainRegistryClient,
	contractRegistryClient contractregistry.ContractRegistryClient,
	identityManagerClient identitymanager.IdentityManagerClient,
	syncProducer sarama.SyncProducer,
	topicCfg *pkgsarama.KafkaTopicConfig,
) (*app.App, error) {
	// Create Data agents
	db, err := multi.Build(context.Background(), cfg.Store, pgmngr)
	if err != nil {
		return nil, err
	}

	ucs := builder.NewUseCases(db, chainRegistryClient, contractRegistryClient, identityManagerClient, syncProducer, topicCfg)

	// Option for transaction handler
	txSchedulerHandlerOpt := app.HandlerOpt(reflect.TypeOf(&dynamic.Transactions{}), controllers.NewBuilder(ucs))

	// Create app
	return app.New(
		cfg.App,
		app.MultiTenancyOpt("auth", jwt, key, cfg.Multitenancy),
		ReadinessOpt(db, chainRegistryClient, identityManagerClient),
		app.MetricsOpt(),
		app.LoggerMiddlewareOpt("base"),
		app.SwaggerOpt("./public/swagger-specs/services/transaction-scheduler/swagger.json", "base@logger-base"),
		txSchedulerHandlerOpt,
		app.ProviderOpt(NewProvider()),
	)
}

func ReadinessOpt(db database.DB, chainRegistryClient chainregistry.ChainRegistryClient, identitymanagerClient identitymanager.IdentityManagerClient) app.Option {
	return func(ap *app.App) error {
		ap.AddReadinessCheck("database", postgres.Checker(db.(orm.DB)))
		ap.AddReadinessCheck("chain-registry", chainRegistryClient.Checker())
		ap.AddReadinessCheck("contract-registry", contractregistry2.GlobalChecker())
		ap.AddReadinessCheck("identity-manager", identitymanagerClient.Checker())
		ap.AddReadinessCheck("kafka", pkgsarama.GlobalClientChecker())
		return nil
	}
}
