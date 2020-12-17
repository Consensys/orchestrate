package txsigner

import (
	"context"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/app"
	pkgsarama "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"
	orchestrateclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	keymanagerclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/nonce"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-signer/service"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-signer/tx-signer/builder"
)

type txSignerDaemon struct {
	keyManagerClient keymanagerclient.KeyManagerClient
	client           orchestrateclient.OrchestrateClient
	ec               ethclient.MultiClient
	nonceManager     nonce.Sender
	consumerGroup    sarama.ConsumerGroup
	producer         sarama.SyncProducer
	config           *Config
}

func NewTxSigner(
	config *Config,
	consumerGroup sarama.ConsumerGroup,
	producer sarama.SyncProducer,
	keyManagerClient keymanagerclient.KeyManagerClient,
	client orchestrateclient.OrchestrateClient,
	ec ethclient.MultiClient,
	nonceManager nonce.Sender,
) (*app.App, error) {
	appli, err := app.New(config.App, readinessOpt(client), app.MetricsOpt())
	if err != nil {
		return nil, err
	}

	txSignerDaemon := &txSignerDaemon{
		keyManagerClient: keyManagerClient,
		client:           client,
		consumerGroup:    consumerGroup,
		producer:         producer,
		config:           config,
		ec:               ec,
		nonceManager:     nonceManager,
	}

	appli.RegisterDaemon(txSignerDaemon)

	return appli, nil
}

func (d *txSignerDaemon) Run(ctx context.Context) error {
	logger := log.WithContext(ctx)
	logger.Infof("starting transaction signer")

	// Create business layer use cases
	useCases := builder.NewUseCases(d.client, d.keyManagerClient, d.ec, d.nonceManager,
		d.config.ChainRegistryURL, d.config.CheckerMaxRecovery)

	// Create service layer listener
	listener := service.NewMessageListener(useCases, d.client, d.producer, d.config.RecoverTopic, d.config.CrafterTopic,
		d.config.BckOff)

	return d.consumerGroup.Consume(ctx, []string{d.config.ListenerTopic}, listener)
}

func (d *txSignerDaemon) Close() error {
	return d.consumerGroup.Close()
}

func readinessOpt(client orchestrateclient.OrchestrateClient) app.Option {
	return func(ap *app.App) error {
		ap.AddReadinessCheck("kafka", pkgsarama.GlobalClientChecker())
		ap.AddReadinessCheck("api", client.Checker())
		return nil
	}
}
