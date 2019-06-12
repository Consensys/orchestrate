package steps

import (
	"context"
	"sync"

	broker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/envelope-store.git/store/grpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/ethclient"

)

// Init inialize handlers
func Init(ctx context.Context) {
	wg := sync.WaitGroup{}

	// Initialize Kafka producer
	wg.Add(1)
	go func() {
		broker.InitSyncProducer(ctx)
		wg.Done()
	}()

	// Initialize Context store
	wg.Add(1)
	go func() {
		grpc.Init(ctx)
		wg.Done()
	}()

	// Initialize Ethereum client
	wg.Add(1)
	go func() {
		ethclient.Init(ctx)
		wg.Done()
	}()

	// Wait for all handlers to be ready
	wg.Wait()
}
