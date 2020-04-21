package configwatcher

import (
	"context"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/chain-registry/use-cases"
)

type Watcher configwatcher.Watcher

func NewWatcher(getChains usecases.GetChains, cfg Config, httpEps *http.EntryPoints) Watcher {
	// Create Configuration Watcher
	// Create configuration listener switching HTTP Endpoints configuration
	listeners := []func(context.Context, interface{}) error{
		httpEps.Switch,
	}

	return configwatcher.New(
		cfg.watcher,
		NewProvider(getChains, cfg.dynamic, time.Second),
		dynamic.Merge,
		listeners,
	)
}
