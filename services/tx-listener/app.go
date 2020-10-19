package txlistener

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	pkgsarama "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	txscheduler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
)

func New(
	cfg *app.Config,
	listener,
	sentry app.Daemon,
) (*app.App, error) {
	appli, err := app.New(
		cfg,
		ReadinessOpt(txscheduler.GlobalClient(), chainregistry.GlobalClient()),
		app.MetricsOpt(),
	)
	if err != nil {
		return nil, err
	}

	appli.RegisterDaemon(listener)
	appli.RegisterDaemon(sentry)

	return appli, nil
}

func ReadinessOpt(txSchedulerClient txscheduler.TransactionSchedulerClient, chainRegistryClient chainregistry.ChainRegistryClient) app.Option {
	return func(ap *app.App) error {
		ap.AddReadinessCheck("chain-registry", chainRegistryClient.Checker())
		ap.AddReadinessCheck("transaction-scheduler", txSchedulerClient.Checker())
		ap.AddReadinessCheck("kafka", pkgsarama.GlobalClientChecker())
		return nil
	}
}
