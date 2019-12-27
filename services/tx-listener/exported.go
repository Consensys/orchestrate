package txlistener

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient/rpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/dynamic"
	provider "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/providers/listener-v1"
	kafkahook "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/session/ethereum/hooks/kafka"
	memoryoffset "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/session/ethereum/offset/memory"
)

var (
	listener *TxListener
	initOnce = &sync.Once{}
)

// TODO: NullProvider should be replaced with chain-registry provider
type NullProvider struct{}

func (p *NullProvider) Run(ctx context.Context, configInput chan<- *dynamic.Message) error {
	return nil
}

// Init hook
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if listener != nil {
			return
		}

		kafkahook.Init(ctx)
		memoryoffset.Init(ctx)
		rpc.Init(ctx)
		provider.Init(ctx)

		listener = NewTxListener(
			provider.GlobalProvider(),
			kafkahook.GlobalHook(),
			memoryoffset.GlobalManager(),
			rpc.GlobalClientV2(),
		)
	})
}

// SetGlobalListener set global TxListener
func SetGlobalListener(l *TxListener) {
	listener = l
}

// GlobalListener return global TxListener
func GlobalListener() *TxListener {
	return listener
}

// Start global TxListener
func Start(ctx context.Context) {
	listener.Start(ctx)
}
