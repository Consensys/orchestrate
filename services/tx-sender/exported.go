package txsender

import (
	"context"

	sarama2 "github.com/Shopify/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/app"
	dbredis "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/redis"
	ethclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient/rpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"
	orchestrateClient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"

	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	keymanager "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/client"
)

// New Utility function used to initialize a new service
func New(ctx context.Context) (*app.App, error) {
	logger := log.FromContext(ctx)
	config := NewConfig(viper.GetViper())

	sarama.InitSyncProducer(ctx)
	keymanager.Init()
	orchestrateClient.Init()
	ethclient.Init(ctx)

	if config.NonceManagerType == NonceManagerTypeRedis {
		dbredis.Init()
	}

	var err error
	consumerGroups := make([]sarama2.ConsumerGroup, config.NConsumer)
	hostnames := viper.GetStringSlice(sarama.KafkaURLViperKey)
	for idx := 0; idx < config.NConsumer; idx++ {
		consumerGroups[idx], err = NewSaramaConsumer(hostnames, config.GroupName)
		if err != nil {
			return nil, err
		}
		logger.WithField("host", hostnames).WithField("group_name", config.GroupName).
			Info("consumer client ready")
	}

	return NewTxSender(
		config,
		consumerGroups,
		sarama.GlobalSyncProducer(),
		keymanager.GlobalClient(),
		orchestrateClient.GlobalClient(),
		ethclient.GlobalClient(),
		dbredis.GlobalClient(),
	)
}

func NewSaramaConsumer(hostnames []string, groupName string) (sarama2.ConsumerGroup, error) {
	config, err := sarama.NewSaramaConfig()
	if err != nil {
		return nil, err
	}

	client, err := sarama.NewClient(hostnames, config)
	if err != nil {
		return nil, err
	}

	return sarama.NewConsumerGroupFromClient(groupName, client)
}
