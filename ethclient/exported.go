package ethclient

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/ethclient/rpc"
)

var (
	client   Client
	initOnce = &sync.Once{}
)

func Init(ctx context.Context) {
	initOnce.Do(func() {
		if client != nil {
			return
		}

		rpc.Init(ctx)

		client = rpc.GlobalClient()
	})
}

// GlobalClient returns global Client
func GlobalClient() Client {
	return client
}

// SetGlobalClient sets global Client
func SetGlobalMultiClient(ec Client) {
	client = ec
}
