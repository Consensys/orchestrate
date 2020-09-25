package txsigner

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
)

func New(
	config *app.Config,
	consumer app.Daemon,
) (*app.App, error) {
	appli, err := app.New(
		config,
		app.MetricsOpt(),
	)

	if err != nil {
		return nil, err
	}

	appli.RegisterDaemon(consumer)

	return appli, nil
}
