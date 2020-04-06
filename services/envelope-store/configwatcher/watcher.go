package configwatcher

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
)

type Watcher configwatcher.Watcher

func NewWatcher(cfg Config, httpEps *http.EntryPoints) Watcher {
	// Create Configuration Watcher
	// Create configuration listener switching HTTP Endpoints configuration
	listeners := []func(context.Context, interface{}) error{
		httpEps.Switch,
	}
	
	return configwatcher.New(
		cfg.watcher,
		NewProvider(cfg.dynamic),
		dynamic.Merge,
		listeners,
	)
}