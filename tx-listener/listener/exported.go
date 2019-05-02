package listener

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/ethclient"
)

var (
	l        TxListener
	conf     *Config
	initOnce = &sync.Once{}
)

func Init(ctx context.Context) {
	initOnce.Do(func() {
		if l != nil {
			return
		}

		// Initialize Ethereum client
		ethclient.Init(context.Background())

		// Initialize Listener
		if conf == nil {
			conf = NewConfig()
		}

		l = NewListener(ethclient.GlobalClient(), conf)

		// Start listening all chains
		for _, chainID := range ethclient.GlobalClient().Networks(context.Background()) {
			_, _ = l.Listen(chainID, -1, 0)
		}
	})
}

func SetGlobalConfig(cfg *Config) {
	conf = cfg
}

// GlobalListener returns global Listener
func GlobalListener() TxListener {
	return l
}

// SetGlobalListener sets global Listener
func SetGlobalListener(listener TxListener) {
	l = listener
}
