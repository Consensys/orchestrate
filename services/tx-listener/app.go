package txlistener

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
)

func New(
	cfg *app.Config,
	listener,
	sentry app.Daemon,
) (*app.App, error) {
	appli, err := app.New(
		cfg,
		app.MetricsOpt(),
	)
	if err != nil {
		return nil, err
	}

	appli.RegisterDaemon(listener)
	appli.RegisterDaemon(sentry)

	return appli, nil
}
