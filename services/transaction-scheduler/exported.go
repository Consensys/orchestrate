package transactionscheduler

import (
	"context"

	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/app"
	authjwt "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/auth/jwt"
	authkey "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/auth/key"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres"
	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/client"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/contract-registry/client"
	identitymanager "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/identity-manager/client"
)

// New Utility function used to initialize a new service
func New(ctx context.Context) (*app.App, error) {
	// Initialize dependencies
	authjwt.Init(ctx)
	authkey.Init(ctx)
	chainregistry.Init(ctx)
	identitymanager.Init()
	sarama.InitSyncProducer(ctx)
	contractregistry.Init(ctx)

	config := NewConfig(viper.GetViper())
	pgmngr := postgres.GetManager()

	return NewTxScheduler(
		config,
		pgmngr,
		authjwt.GlobalChecker(),
		authkey.GlobalChecker(),
		chainregistry.GlobalClient(),
		contractregistry.GlobalClient(),
		identitymanager.GlobalClient(),
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
