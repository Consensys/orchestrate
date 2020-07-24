package contractregistry

import (
	"context"
	"sync"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	authjwt "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/jwt"
	authkey "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/key"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
)

var (
	initOnce = &sync.Once{}
	appli    *app.App
)

func Init(ctx context.Context) {
	initOnce.Do(func() {
		// Initialize dependencies
		authjwt.Init(ctx)
		authkey.Init(ctx)

		var err error
		appli, err = New(
			NewConfig(viper.GetViper()),
			postgres.GetManager(),
			authjwt.GlobalChecker(), authkey.GlobalChecker(),
		)
		if err != nil {
			log.FromContext(ctx).WithError(err).Fatalf("Could not create contract registry application")
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
