package txsender

import (
	"context"

	sarama2 "github.com/Shopify/sarama"
	orchestrateClient "github.com/consensys/orchestrate/pkg/sdk/client"
	"github.com/consensys/orchestrate/pkg/toolkit/app"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	dbredis "github.com/consensys/orchestrate/pkg/toolkit/database/redis"
	ethclient "github.com/consensys/orchestrate/pkg/toolkit/ethclient/rpc"

	"github.com/consensys/orchestrate/pkg/broker/sarama"
	qkm "github.com/consensys/orchestrate/pkg/quorum-key-manager"
	"github.com/spf13/viper"
)

// New Utility function used to initialize a new service
func New(ctx context.Context) (*app.App, error) {
	logger := log.FromContext(ctx)
	config := NewConfig(viper.GetViper())

	sarama.InitSyncProducer(ctx)
	qkm.Init()
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
		qkm.GlobalClient(),
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
