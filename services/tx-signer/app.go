package txsigner

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v4"
	pkgsarama "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	client2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-signer/tx-signer/builder"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-signer/service"
)

type txSignerDaemon struct {
	keyManagerClient  client.KeyManagerClient
	txSchedulerClient client2.TransactionSchedulerClient
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
	}
	appli.RegisterDaemon(txSignerDaemon)

	return appli, nil
}

func (signer *txSignerDaemon) Run(ctx context.Context) error {
	logger := log.WithContext(ctx)
	logger.Infof("starting transaction signer")

	// Create business layer use cases
	useCases := builder.NewUseCases(signer.keyManagerClient, signer.producer)

	// Create service layer listener
	listener := service.NewMessageListener(useCases, signer.config.SenderTopic, signer.config.RecoverTopic, signer.txSchedulerClient, retryMessageBackOff())

	return signer.consumerGroup.Consume(ctx, []string{signer.config.ListenerTopic}, listener)
}

func (signer *txSignerDaemon) Close() error {
	return signer.consumerGroup.Close()
}

func readinessOpt(txSchedulerClient client2.TransactionSchedulerClient) app.Option {
	return func(ap *app.App) error {
		ap.AddReadinessCheck("kafka", pkgsarama.GlobalClientChecker())
		ap.AddReadinessCheck("transaction-scheduler", txSchedulerClient.Checker())
		return nil
	}
}

func retryMessageBackOff() backoff.BackOff {
	bckOff := backoff.NewExponentialBackOff()
	bckOff.MaxInterval = time.Second * 15
	bckOff.MaxElapsedTime = time.Minute * 5
	return bckOff
}
