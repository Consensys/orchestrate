package api

import (
	"context"

	ethclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient/rpc"

	keymanager "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/client"

	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/app"
	authjwt "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/auth/jwt"
	authkey "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/auth/key"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres"
)

// New Utility function used to initialize a new service
func New(ctx context.Context) (*app.App, error) {
	// Initialize dependencies
	authjwt.Init(ctx)
	authkey.Init(ctx)
	sarama.InitSyncProducer(ctx)
	keymanager.Init()
	ethclient.Init(ctx)

	config := NewConfig(viper.GetViper())
	pgmngr := postgres.GetManager()

	return NewAPI(
		config,
		pgmngr,
		authjwt.GlobalChecker(),
		authkey.GlobalChecker(),
		keymanager.GlobalClient(),
		ethclient.GlobalClient(),
		sarama.GlobalSyncProducer(),
		sarama.NewKafkaTopicConfig(viper.GetViper()),
	)
}

func Run(ctx context.Context) error {
	appli, err := New(ctx)
	if err != nil {
		return err
	}
	return appli.Run(ctx)
}
