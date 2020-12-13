package txsigner

import (
	"context"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/app"
	pkgsarama "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/nonce"
	client2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-signer/service"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-signer/tx-signer/builder"
)

type txSignerDaemon struct {
	keyManagerClient  client.KeyManagerClient
	txSchedulerClient client2.TransactionSchedulerClient
	ec                ethclient.MultiClient
	nonceManager      nonce.Sender
	consumerGroup     sarama.ConsumerGroup
	producer          sarama.SyncProducer
	config            *Config
}

func NewTxSigner(
	config *Config,
	consumerGroup sarama.ConsumerGroup,
	producer sarama.SyncProducer,
	keyManagerClient client.KeyManagerClient,
	txSchedulerClient client2.TransactionSchedulerClient,
	ec ethclient.MultiClient,
	nonceManager nonce.Sender,
) (*app.App, error) {
	appli, err := app.New(config.App, readinessOpt(txSchedulerClient), app.MetricsOpt())
	if err != nil {
		return nil, err
	}

	txSignerDaemon := &txSignerDaemon{
		keyManagerClient:  keyManagerClient,
		txSchedulerClient: txSchedulerClient,
		consumerGroup:     consumerGroup,
		producer:          producer,
		config:            config,
		ec:                ec,
		nonceManager:      nonceManager,
	}

	appli.RegisterDaemon(txSignerDaemon)

	return appli, nil
}

func (d *txSignerDaemon) Run(ctx context.Context) error {
	logger := log.WithContext(ctx)
	logger.Infof("starting transaction signer")

	// Create business layer use cases
	useCases := builder.NewUseCases(d.txSchedulerClient, d.keyManagerClient, d.ec, d.nonceManager,
		d.config.ChainRegistryURL, d.config.CheckerMaxRecovery)

	// Create service layer listener
	listener := service.NewMessageListener(useCases, d.txSchedulerClient, d.producer, d.config.RecoverTopic, d.config.CrafterTopic,
		d.config.BckOff)

	return d.consumerGroup.Consume(ctx, []string{d.config.ListenerTopic}, listener)
}

func (d *txSignerDaemon) Close() error {
	return d.consumerGroup.Close()
}

func readinessOpt(txSchedulerClient client2.TransactionSchedulerClient) app.Option {
	return func(ap *app.App) error {
		ap.AddReadinessCheck("kafka", pkgsarama.GlobalClientChecker())
		ap.AddReadinessCheck("transaction-scheduler", txSchedulerClient.Checker())
		return nil
	}
}
