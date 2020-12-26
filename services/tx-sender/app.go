package txsender

import (
	"context"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/app"
	pkgsarama "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	dbredis "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/redis"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"
	api "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	keymanager "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-sender/service"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-sender/store/memory"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-sender/store/redis"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-sender/tx-sender/builder"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-sender/tx-sender/nonce"
)

type txSenderDaemon struct {
	keyManagerClient keymanager.KeyManagerClient
	jobClient        api.JobClient
	ec               ethclient.MultiClient
	nonceManager     nonce.Manager
	consumerGroup    sarama.ConsumerGroup
	producer         sarama.SyncProducer
	config           *Config
}

func NewTxSender(
	config *Config,
	consumerGroup sarama.ConsumerGroup,
	producer sarama.SyncProducer,
	keyManagerClient keymanager.KeyManagerClient,
	apiClient api.OrchestrateClient,
	ec ethclient.MultiClient,
	redisCli *dbredis.Client,
) (*app.App, error) {
	appli, err := app.New(config.App, readinessOpt(apiClient, redisCli), app.MetricsOpt())
	if err != nil {
		return nil, err
	}

	var nm nonce.Manager
	if config.NonceManagerType == NonceManagerTypeInMemory {
		nm = nonce.NewNonceManager(ec, memory.NewNonceSender(), memory.NewNonceRecoveryTracker(),
			config.ChainRegistryURL, config.NonceMaxRecovery)
	} else if config.NonceManagerType == NonceManagerTypeRedis {
		nm = nonce.NewNonceManager(ec, redis.NewNonceSender(redisCli), redis.NewNonceRecoveryTracker(redisCli),
			config.ChainRegistryURL, config.NonceMaxRecovery)
	}

	txSenderDaemon := &txSenderDaemon{
		keyManagerClient: keyManagerClient,
		jobClient:        apiClient,
		consumerGroup:    consumerGroup,
		producer:         producer,
		config:           config,
		ec:               ec,
		nonceManager:     nm,
	}

	appli.RegisterDaemon(txSenderDaemon)

	return appli, nil
}

func (d *txSenderDaemon) Run(ctx context.Context) error {
	logger := log.WithContext(ctx)
	logger.Infof("starting transaction sender")

	// Create business layer use cases
	useCases := builder.NewUseCases(d.jobClient, d.keyManagerClient, d.ec, d.nonceManager,
		d.config.ChainRegistryURL, d.config.NonceMaxRecovery)

	// Create service layer listener
	listener := service.NewMessageListener(useCases, d.jobClient, d.producer, d.config.RecoverTopic, d.config.SenderTopic,
		d.config.BckOff)

	return d.consumerGroup.Consume(ctx, []string{d.config.SenderTopic}, listener)
}

func (d *txSenderDaemon) Close() error {
	return d.consumerGroup.Close()
}

func readinessOpt(apiClient api.MetricClient, redisCli *dbredis.Client) app.Option {
	return func(ap *app.App) error {
		ap.AddReadinessCheck("kafka", pkgsarama.GlobalClientChecker())
		ap.AddReadinessCheck("api", apiClient.Checker())
		if redisCli != nil {
			ap.AddReadinessCheck("redis", redisCli.Ping)
		}
		return nil
	}
}
