package chainregistry

import (
	"time"

	prom "github.com/prometheus/client_golang/prometheus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient"
	metrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/multi"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/chain-registry/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/service/configwatcher"
	ctrl "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/service/controllers"
	chainctrl "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/service/controllers/chains"
	faucetctrl "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/service/controllers/faucets"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/service/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
)

func newApplication(
	cfg *Config,
	dataAgents store.DataAgents,
	ethClient ethclient.ChainLedgerReader,
	jwt, key auth.Checker,
	promRegisterer prom.Registerer,
) (*app.App, error) {
	// Create metrics registry and register it on prometheus
	reg := metrics.New(cfg.app.Metrics)
	err := promRegisterer.Register(reg.Prometheus())

	if err != nil {
		return nil, err
	}

	getChainsUC := usecases.NewGetChains(dataAgents.Chain)

	// Create HTTP Handler for Chain
	chainCtrl := chainctrl.NewController(
		getChainsUC,
		usecases.NewGetChain(dataAgents.Chain),
		usecases.NewRegisterChain(dataAgents.Chain, ethClient),
		usecases.NewDeleteChain(dataAgents.Chain),
		usecases.NewUpdateChain(dataAgents.Chain),
	)

	// Create HTTP Handler for Faucet
	faucetCtrl := faucetctrl.NewController(
		usecases.NewGetFaucets(dataAgents.Faucet),
		usecases.NewGetFaucet(dataAgents.Faucet),
		usecases.NewRegisterFaucet(dataAgents.Faucet),
		usecases.NewDeleteFaucet(dataAgents.Faucet),
		usecases.NewUpdateFaucet(dataAgents.Faucet),
	)

	builderCtrl := ctrl.NewBuilder(chainCtrl, faucetCtrl)
	routerBuilder, err := http.NewHTTPRouterBuilder(builderCtrl, cfg.app.HTTP, jwt, key, cfg.multitenancy, reg.HTTP())
	if err != nil {
		return nil, err
	}

	// Create app
	return app.New(
		cfg.app,
		configwatcher.NewProvider(getChainsUC, configwatcher.NewInternalConfig(cfg.app.HTTP, cfg.app.Watcher).DynamicCfg(), time.Second),
		routerBuilder,
		nil,
		reg,
	)
}
