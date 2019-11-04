package listener

import (
	"context"
	"math/big"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/tx-listener/handler"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/tx-listener/listener/base"
)

var (
	l        TxListener
	conf     *base.Config
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
			conf = base.NewConfig()
		}

		l = base.NewTxListener(ethclient.GlobalClient(), conf)

		go func() {
			<-ctx.Done()
			l.Close()
		}()
	})
}

func SetGlobalConfig(cfg *base.Config) {
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

// Listen start listening
func Listen(ctx context.Context, chains []*big.Int, h handler.TxListenerHandler) error {
	return l.Listen(ctx, chains, h)
}
