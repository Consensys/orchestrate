package sender

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	contextStore "gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/store"
	storegrpc "gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/store/grpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	types "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/context-store"
)

var (
	store contextStore.EnvelopeStore

	handler  engine.HandlerFunc
	initOnce = &sync.Once{}
)

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
	store = storegrpc.NewEnvelopeStore(types.NewStoreClient(conn))
	log.WithFields(log.Fields{
		"grpc.store.target": conn.Target(),
	}).Infof("infra-store: grpc client connected")

	// TODO: properly close connection
	// go func() {
	// 	// Close connection when infrastructure closes
	// 	<-infra.ctx.Done()
	// 	conn.Close()
	// }()
}

// Init initialize Gas Pricer Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if handler != nil {
			return
		}

		// Initialize Context store
		initStore(ctx)

		// Initialize Ethereum client
		ethclient.Init(ctx)

		// Create Handler
		handler = Sender(ethclient.GlobalMultiClient(), store)

		log.Infof("sender: handler ready")
	})
}

// SetGlobalHandler sets global Gas Pricer Handler
func SetGlobalHandler(h engine.HandlerFunc) {
	handler = h
}

// GlobalHandler returns global Gas Pricer handler
func GlobalHandler() engine.HandlerFunc {
	return handler
}
