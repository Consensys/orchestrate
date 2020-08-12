package transactionscheduler

import (
	"context"

	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/client"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"

	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	authjwt "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/jwt"
	authkey "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/key"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
)

type TxScheduler struct {
	TxSchedulerAPI app.Service
	TxSentryDaemon app.Daemon
}

func NewTxScheduler(ctx context.Context) (*TxScheduler, error) {
	// Initialize dependencies
	authjwt.Init(ctx)
	authkey.Init(ctx)
	client.Init(ctx)
	sarama.InitSyncProducer(ctx)
	contractregistry.Init(ctx, viper.GetString(contractregistry.ContractRegistryURLViperKey))

	config := NewConfig(viper.GetViper())
	pgmngr := postgres.GetManager()

	txSchedulerApp, err := NewTxSchedulerApp(
		config,
		pgmngr,
		authjwt.GlobalChecker(), authkey.GlobalChecker(),
		client.GlobalClient(),
		contractregistry.GlobalClient(),
		sarama.GlobalSyncProducer(),
		sarama.NewKafkaTopicConfig(viper.GetViper()),
	)
	if err != nil {
		return nil, err
	}

	txSentryDaemon, err := NewTxSentryDaemon(pgmngr, config)
	if err != nil {
		return nil, err
	}

	return &TxScheduler{
		TxSchedulerAPI: txSchedulerApp,
		TxSentryDaemon: txSentryDaemon,
	}, nil
}

func (txScheduler *TxScheduler) StartSentry(ctx context.Context) chan error {
	return txScheduler.TxSentryDaemon.Start(ctx)
}

func (txScheduler *TxScheduler) StartScheduler(ctx context.Context) error {
	return txScheduler.TxSchedulerAPI.Start(ctx)
}

func (txScheduler *TxScheduler) StopSentry(ctx context.Context) {
	txScheduler.TxSentryDaemon.Stop(ctx)
}

func (txScheduler *TxScheduler) StopScheduler(ctx context.Context) error {
	return txScheduler.TxSchedulerAPI.Stop(ctx)
}

func (txScheduler *TxScheduler) IsReady() bool {
	return txScheduler.TxSentryDaemon.IsReady() && txScheduler.TxSchedulerAPI.IsReady()
}
