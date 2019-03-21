package infra

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	storegrpc "gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/infra/grpc"
	store "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/context-store"
	"google.golang.org/grpc"
)

func initStore(infra *Infra) {
	// Init grpc store
	conn, err := grpc.DialContext(
		context.Background(),
		viper.GetString("grpc.store.target"),
		grpc.WithInsecure(),
	)

	if err != nil {
		log.WithError(err).Fatalf("infra-store: failed to dial grpc server")
	}

	// Set store
	infra.Store = storegrpc.NewTraceStore(store.NewStoreClient(conn))
	log.WithFields(log.Fields{
		"grpc.store.target": conn.Target(),
	}).Infof("infra-store: grpc client connected")

	go func() {
		// Close connection when infrastructure closes
		<-infra.ctx.Done()
		conn.Close()
	}()
}
