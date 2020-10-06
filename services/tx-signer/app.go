package txsigner

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	pkgsarama "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/secretstore"
)

func New(
	config *app.Config,
	consumer app.Daemon,
) (*app.App, error) {
	appli, err := app.New(
		config,
		ReadinessOpt(),
		app.MetricsOpt(),
	)

	if err != nil {
		return nil, err
	}

	appli.RegisterDaemon(consumer)

	return appli, nil
}

func ReadinessOpt() app.Option {
	return func(ap *app.App) error {
		ap.AddReadinessCheck("kafka", pkgsarama.GlobalClientChecker())
		ap.AddReadinessCheck("secret-store", secretstore.GlobalSecretStoreChecker())
		return nil
	}
}
