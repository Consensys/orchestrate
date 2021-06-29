package txsender

import (
	"context"
	"fmt"
	"time"

	pkgsarama "github.com/ConsenSys/orchestrate/pkg/broker/sarama"
	"github.com/ConsenSys/orchestrate/pkg/errors"
	keymanager "github.com/ConsenSys/orchestrate/pkg/quorum-key-manager/client"
	api "github.com/ConsenSys/orchestrate/pkg/sdk/client"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
	dbredis "github.com/ConsenSys/orchestrate/pkg/toolkit/database/redis"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/ethclient"
	"github.com/ConsenSys/orchestrate/services/tx-sender/service"
	"github.com/ConsenSys/orchestrate/services/tx-sender/store/memory"
	"github.com/ConsenSys/orchestrate/services/tx-sender/store/redis"
	"github.com/ConsenSys/orchestrate/services/tx-sender/tx-sender/builder"
	"github.com/ConsenSys/orchestrate/services/tx-sender/tx-sender/nonce"
	"github.com/Shopify/sarama"
	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/go-multierror"
)

const component = "application"

type txSenderDaemon struct {
	keyManagerClient keymanager.KeyManagerClient
	jobClient        api.JobClient
	ec               ethclient.MultiClient
	nonceManager     nonce.Manager
	consumerGroup    []sarama.ConsumerGroup
	producer         sarama.SyncProducer
	config           *Config
	logger           *log.Logger
	cancel           context.CancelFunc
}

func NewTxSender(
	config *Config,
	consumerGroup []sarama.ConsumerGroup,
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
	d.logger.Debug("starting transaction sender")

	// Create business layer use cases
	useCases := builder.NewUseCases(d.jobClient, d.keyManagerClient, d.ec, d.nonceManager,
		d.config.ProxyURL, d.config.NonceMaxRecovery)

	// Create service layer listener
	listener := service.NewMessageListener(useCases, d.jobClient, d.producer, d.config.RecoverTopic, d.config.SenderTopic,
		d.config.BckOff)

	ctx, d.cancel = context.WithCancel(ctx)
	gr := &multierror.Group{}
	for idx, consumerGroup := range d.consumerGroup {
		cGroup := consumerGroup
		cGroupID := fmt.Sprintf("c-%d", idx)
		logger := d.logger.WithField("consumer", cGroupID)
		cctx := log.With(log.WithField(ctx, "consumer", cGroupID), logger)
		gr.Go(func() error {
			// We retry once after consume exits to prevent entire stack to exit after kafka rebalance is triggered
			err := backoff.RetryNotify(
				func() error {
					err := cGroup.Consume(cctx, []string{d.config.SenderTopic}, listener)

					// In this case, kafka rebalance was triggered and we want to retry
					if err == nil && cctx.Err() == nil {
						return fmt.Errorf("kafka rebalance was triggered")
					}

					return backoff.Permanent(err)
				},
				backoff.NewConstantBackOff(time.Millisecond*500),
				func(err error, duration time.Duration) {
					logger.WithError(err).Warnf("consuming session exited, retrying in %s", duration.String())
				},
			)
			d.cancel()
			return err
		})
	}

	return gr.Wait().ErrorOrNil()
}

func (d *txSenderDaemon) Close() error {
	var gerr error
	for _, consumerGroup := range d.consumerGroup {
		gerr = errors.CombineErrors(gerr, consumerGroup.Close())
	}

	return gerr
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
