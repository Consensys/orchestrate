package listenerv1

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient/rpc"
)

var (
	provider *Provider
	initOnce = &sync.Once{}
)

// Init hook
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if provider != nil {
			return
		}

		conf, err := NewConfig()
		if err != nil {
			log.WithError(err).Fatalf("tx-listener: invalid configuration")
		}

		rpc.Init(ctx)

		provider = NewProvider(conf, rpc.GlobalClientV2())
	})
}

// SetGlobalProvider set global TxProvider
func SetGlobalProvider(p *Provider) {
	provider = p
}

// GlobalProvider return global TxProvider
func GlobalProvider() *Provider {
	return provider
}
