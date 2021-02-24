package txlistener

import (
	"github.com/ConsenSys/orchestrate/pkg/app"
	pkgsarama "github.com/ConsenSys/orchestrate/pkg/broker/sarama"
	orchestrateclient "github.com/ConsenSys/orchestrate/pkg/sdk/client"
	listenermetrics "github.com/ConsenSys/orchestrate/services/tx-listener/metrics"
	provider "github.com/ConsenSys/orchestrate/services/tx-listener/providers"
	"github.com/ConsenSys/orchestrate/services/tx-listener/session/ethereum"
	hook "github.com/ConsenSys/orchestrate/services/tx-listener/session/ethereum/hooks"
	"github.com/ConsenSys/orchestrate/services/tx-listener/session/ethereum/offset"
	txsentry "github.com/ConsenSys/orchestrate/services/tx-sentry"
	"github.com/spf13/viper"
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
