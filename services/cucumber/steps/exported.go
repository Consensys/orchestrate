package steps

import (
	"context"
	"sync"

	broker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
	grpcStore "gitlab.com/ConsenSys/client/fr/core-stack/service/envelope-store.git/store/grpc"
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
		grpcStore.Init(ctx)
		wg.Done()
	}()

	// Wait for all handlers to be ready
	wg.Wait()
}
