package grpc

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	types "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope-store"
	"google.golang.org/grpc"
)

var (
	envelopeStore *EnvelopeStore
	initOnce      = &sync.Once{}
)

// InitStore initilialize envelope store
func initStore(ctx context.Context) {
	// Init grpc store
	conn, err := grpc.DialContext(
		ctx,
		viper.GetString("grpc.store.target"),
		grpc.WithInsecure(),
	)

	if err != nil {
		log.WithError(err).Fatalf("infra-store: failed to dial grpc server")
	}

	// Set store
	envelopeStore = NewEnvelopeStore(types.NewStoreClient(conn))
	log.WithFields(log.Fields{
		"grpc.store.target": conn.Target(),
	}).Infof("infra-store: grpc client connected")

	go func() {
		// Close connection when infrastructure closes
		<-ctx.Done()
		conn.Close()
	}()
}

// Init initialize Sender Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if envelopeStore != nil {
			return
		}

		// Initialize Grpc store
		initStore(ctx)

		log.Infof("grpc: store ready")
	})
}

func GlobalEnvelopeStore() *EnvelopeStore {
	return envelopeStore
}

// SetGlobalConfig sets Sarama global configuration
func SetGlobalEnvelopeStore(s *EnvelopeStore) {
	envelopeStore = s
}
