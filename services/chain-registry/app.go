package chainregistry

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
)

func New(
	cfg *app.Config,
	jwt, key auth.Checker,
	multitenancy bool,
	s store.ChainRegistryStore,
	ec ethclient.ChainLedgerReader,
) (*app.App, error) {
	// Create HTTP Router builder
	httpBuilder, err := NewHTTPBuilder(cfg.HTTP, jwt, key, multitenancy, s, ec)
	if err != nil {
		return nil, err
	}

	// Create HTTP EntryPoints
	httpEps := http.NewEntryPoints(
		cfg.HTTP.EntryPoints,
		httpBuilder,
	)

	// Create Configuration Watcher
	// Create configuration listener switching HTTP Endpoints configuration
	listeners := []func(context.Context, interface{}) error{
		httpEps.Switch,
	}

	watcher := configwatcher.New(
		cfg.Watcher,
		NewProvider(cfg.HTTP, s),
		dynamic.Merge,
		listeners,
	)

	// Create app
	return app.New(watcher, httpEps, nil), nil
}

func NewStore(pgmngr postgres.Manager, conf *store.Config) (store.ChainRegistryStore, error) {
	storeBuilder := store.NewBuilder(pgmngr)
	return storeBuilder.Build(context.Background(), conf)
}
