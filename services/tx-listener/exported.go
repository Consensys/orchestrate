package txlistener

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient/rpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/dynamic"
	registryprovider "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/providers/chain-registry"
	kafkahook "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/session/ethereum/hooks/kafka"
	registryoffset "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/session/ethereum/offset/chain-registry"
)

var (
	listener *TxListener
	initOnce = &sync.Once{}
)

// TODO: NullProvider should be replaced with chain-registry provider
type NullProvider struct{}

func (p *NullProvider) Run(_ context.Context, _ chan<- *dynamic.Message) error {
	return nil
}

// Init hook
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if listener != nil {
			return
		}

		registryprovider.Init(ctx)
		kafkahook.Init(ctx)
		registryoffset.Init(ctx)
		rpc.Init(ctx)

		listener = NewTxListener(
			registryprovider.GlobalProvider(),
			kafkahook.GlobalHook(),
			registryoffset.GlobalManager(),
			rpc.GlobalClient(),
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
