package transactionscheduler

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"

	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/client"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"

	"github.com/spf13/viper"
	authjwt "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/jwt"
	authkey "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/key"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
)

// New Utility function used to initialize a new service
func New(ctx context.Context) (*app.App, error) {
	// Initialize dependencies
	authjwt.Init(ctx)
	authkey.Init(ctx)
	client.Init(ctx)
	sarama.InitSyncProducer(ctx)
	contractregistry.Init(ctx, viper.GetString(contractregistry.ContractRegistryURLViperKey))

	config := NewConfig(viper.GetViper())
	pgmngr := postgres.GetManager()

	return NewTxScheduler(
		config,
		pgmngr,
		authjwt.GlobalChecker(), authkey.GlobalChecker(),
		client.GlobalClient(),
		contractregistry.GlobalClient(),
		sarama.GlobalSyncProducer(),
		sarama.NewKafkaTopicConfig(viper.GetViper()),
	)
}
