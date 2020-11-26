package txlistener

import (
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/app"
	pkgsarama "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/client"
	txscheduler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/client"
	listenermetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-listener/metrics"
	provider "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-listener/providers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-listener/session/ethereum"
	hook "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-listener/session/ethereum/hooks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-listener/session/ethereum/offset"
	txsentry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-sentry"
)

func New(
	cfg *app.Config,
	prvdr provider.Provider,
	hk hook.Hook,
	offsets offset.Manager,
	ec ethereum.EthClient,
	txSchedulerClientListener,
	txSchedulerClientSentry txscheduler.TransactionSchedulerClient,

) (*app.App, error) {

	var listenerMetrics listenermetrics.ListenerMetrics
	if cfg.Metrics.IsActive(listenermetrics.ModuleName) {
		listenerMetrics = listenermetrics.NewListenerMetrics()
	} else {
		listenerMetrics = listenermetrics.NewListenerNopMetrics()
	}

	listener = NewTxListener(
		prvdr,
		hk,
		offsets,
		ec,
		txSchedulerClientListener,
		listenerMetrics,
	)

	sentry = txsentry.NewTxSentry(
		txSchedulerClientSentry,
		txsentry.NewConfig(viper.GetViper()),
	)

	appli, err := app.New(
		cfg,
		ReadinessOpt(txscheduler.GlobalClient(), chainregistry.GlobalClient()),
		app.MetricsOpt(listenerMetrics),
	)
	if err != nil {
		return nil, err
	}

	appli.RegisterDaemon(listener)
	appli.RegisterDaemon(sentry)

	return appli, nil
}

func ReadinessOpt(txSchedulerClient txscheduler.TransactionSchedulerClient, chainRegistryClient chainregistry.ChainRegistryClient) app.Option {
	return func(ap *app.App) error {
		ap.AddReadinessCheck("chain-registry", chainRegistryClient.Checker())
		ap.AddReadinessCheck("transaction-scheduler", txSchedulerClient.Checker())
		ap.AddReadinessCheck("kafka", pkgsarama.GlobalClientChecker())
		return nil
	}
}
