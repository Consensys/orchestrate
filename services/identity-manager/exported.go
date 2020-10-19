package identitymanager

import (
	"context"

	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	authjwt "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/jwt"
	authkey "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/key"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	chainRegistryClient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	keyManagerClient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/client"
	txSchedulerClient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
)

// New Utility function used to initialize a new service
func New(ctx context.Context) (*app.App, error) {
	// Initialize dependencies
	authjwt.Init(ctx)
	authkey.Init(ctx)
	keyManagerClient.Init()
	chainRegistryClient.Init(ctx)
	txSchedulerClient.Init()

	config := NewConfig(viper.GetViper())
	pgmngr := postgres.GetManager()

	return NewIdentityManager(config, pgmngr, authjwt.GlobalChecker(), authkey.GlobalChecker(),
		keyManagerClient.GlobalClient(), chainRegistryClient.GlobalClient(), txSchedulerClient.GlobalClient())
}

func Run(ctx context.Context) error {
	appli, err := New(ctx)
	if err != nil {
		return err
	}
	return appli.Run(ctx)
}
