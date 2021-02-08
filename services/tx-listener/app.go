package txlistener

import (
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/app"
	pkgsarama "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	orchestrateclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
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
	client orchestrateclient.OrchestrateClient,
) (*app.App, error) {

	var listenerMetrics listenermetrics.ListenerMetrics
	if cfg.Metrics.IsActive(listenermetrics.ModuleName) {
		listenerMetrics = listenermetrics.NewListenerMetrics()
	} else {
		listenerMetrics = listenermetrics.NewListenerNopMetrics()
	}

	listener = NewTxListener(prvdr, hk, offsets, ec, client, listenerMetrics)
	sentry = txsentry.NewTxSentry(client, txsentry.NewConfig(viper.GetViper()))
	appli, err := app.New(cfg, ReadinessOpt(client), app.MetricsOpt(listenerMetrics))
	if err != nil {
		return nil, err
	}

	appli.RegisterDaemon(listener)
	appli.RegisterDaemon(sentry)

	return appli, nil
}

func ReadinessOpt(client orchestrateclient.OrchestrateClient) app.Option {
	return func(ap *app.App) error {
		ap.AddReadinessCheck("api", client.Checker())
		ap.AddReadinessCheck("kafka", pkgsarama.GlobalClientChecker())
		return nil
	}
}
