package txsender

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/app"
	dbredis "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/redis"
	ethclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient/rpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"

	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	keymanager "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/client"
)

// New Utility function used to initialize a new service
func New(ctx context.Context) (*app.App, error) {
	config := NewConfig(viper.GetViper())

	// Initialize dependencies
	sarama.InitConsumerGroup(ctx, config.GroupName)
	sarama.InitSyncProducer(ctx)
	keymanager.Init()
	client.Init()
	ethclient.Init(ctx)

	if config.NonceManagerType == NonceManagerTypeRedis {
		dbredis.Init()
	}

	return NewTxSender(
		config,
		sarama.GlobalConsumerGroup(),
		sarama.GlobalSyncProducer(),
		keymanager.GlobalClient(),
		client.GlobalClient(),
		ethclient.GlobalClient(),
		dbredis.GlobalClient(),
	)
}
