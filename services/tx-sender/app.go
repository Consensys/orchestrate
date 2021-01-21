package txsender

import (
	"context"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/cenkalti/backoff/v4"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/app"
	pkgsarama "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	dbredis "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/redis"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"
	api "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	keymanager "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-sender/service"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-sender/store/memory"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-sender/store/redis"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-sender/tx-sender/builder"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-sender/tx-sender/nonce"
)

const component = "application"

type txSenderDaemon struct {
	keyManagerClient keymanager.KeyManagerClient
	jobClient        api.JobClient
	ec               ethclient.MultiClient
	nonceManager     nonce.Manager
	consumerGroup    sarama.ConsumerGroup
	producer         sarama.SyncProducer
	config           *Config
	logger           *log.Logger
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
		nm = nonce.NewNonceManager(ec, memory.NewNonceSender(config.NonceManagerExpiration), memory.NewNonceRecoveryTracker(),
			config.ProxyURL, config.NonceMaxRecovery)
	} else if config.NonceManagerType == NonceManagerTypeRedis {
		nm = nonce.NewNonceManager(ec, redis.NewNonceSender(redisCli), redis.NewNonceRecoveryTracker(redisCli),
			config.ProxyURL, config.NonceMaxRecovery)
	}

	txSenderDaemon := &txSenderDaemon{
		keyManagerClient: keyManagerClient,
		jobClient:        apiClient,
		consumerGroup:    consumerGroup,
		producer:         producer,
		config:           config,
		ec:               ec,
		nonceManager:     nm,
		logger:           log.NewLogger().SetComponent(component),
	}

	appli.RegisterDaemon(txSenderDaemon)

	return appli, nil
}

func (d *txSenderDaemon) Run(ctx context.Context) error {
	logger := d.logger.WithContext(ctx)
	logger.Debug("starting transaction sender")

	// Create business layer use cases
	useCases := builder.NewUseCases(d.jobClient, d.keyManagerClient, d.ec, d.nonceManager,
		d.config.ProxyURL, d.config.NonceMaxRecovery)

	// Create service layer listener
	listener := service.NewMessageListener(useCases, d.jobClient, d.producer, d.config.RecoverTopic, d.config.SenderTopic,
		d.config.BckOff)

	// We retry once after consume exits to prevent entire stack to exit after kafka rebalance is triggered
	return backoff.RetryNotify(
		func() error {
			err := d.consumerGroup.Consume(ctx, []string{d.config.SenderTopic}, listener)

			// In this case, kafka rebalance was triggered and we want to retry
			if err == nil && ctx.Err() == nil {
				return fmt.Errorf("kafka rebalance was triggered")
			}

			return backoff.Permanent(err)
		},
		backoff.NewConstantBackOff(time.Millisecond*500),
		func(err error, duration time.Duration) {
			logger.WithError(err).Warnf("consuming session exited, retrying in %s", duration.String())
		},
	)
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
