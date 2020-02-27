package chains

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient/rpc"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
)

var (
	handler  *Handler
	initOnce = &sync.Once{}
)

// Initialize API handlers
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if handler != nil {
			return
		}

		store.Init(ctx)
		rpc.Init(ctx)

		// Set Chain-Registry handler
		handler = NewHandler(store.GlobalStoreRegistry(), rpc.GlobalClient())
	})
}

// GlobalChainRegistryClient return the chain registry
func GlobalHandler() *Handler {
	return handler
}

// SetGlobalChainRegistryClient set a the chain registry client
func SetGlobalHandler(h *Handler) {
	handler = h
}
