package chainregistry

import (
	"sync"

	orchestrateclient "github.com/ConsenSys/orchestrate/pkg/sdk/client"
)

var (
	provider *Provider
	initOnce = &sync.Once{}
)

// Init hook
func Init(client orchestrateclient.OrchestrateClient) {
	initOnce.Do(func() {
		if provider != nil {
			return
		}

		provider = &Provider{
			conf:   NewConfig(),
			client: client,
		}
	})
}

// GlobalSetProvider returns global a chain-registry provider
func GlobalProvider() *Provider {
	return provider
}
