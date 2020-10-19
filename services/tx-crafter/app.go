package txcrafter

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	pkgsarama "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/redis"
	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	txscheduler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
)

func New(
	config *app.Config,
	consumer app.Daemon,
) (*app.App, error) {
	appli, err := app.New(
		config,
		ReadinessOpt(txscheduler.GlobalClient(), chainregistry.GlobalClient()),
		app.MetricsOpt(),
	)

	if err != nil {
		return nil, err
	}

	appli.RegisterDaemon(consumer)

	return appli, nil
}

func ReadinessOpt(txSchedulerClient txscheduler.TransactionSchedulerClient, chainRegistryClient chainregistry.ChainRegistryClient) app.Option {
	return func(ap *app.App) error {
		ap.AddReadinessCheck("chain-registry", chainRegistryClient.Checker())
		ap.AddReadinessCheck("transaction-scheduler", txSchedulerClient.Checker())
		ap.AddReadinessCheck("kafka", pkgsarama.GlobalClientChecker())
		ap.AddReadinessCheck("redis", redis.GlobalChecker())
		return nil
	}
}
