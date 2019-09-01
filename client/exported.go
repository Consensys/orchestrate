package client

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
	grpcclient "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/grpc/client"
	evlpstore "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services/envelope-store"
	"google.golang.org/grpc"
)

const component = "envelope-store.client"

var (
	client   evlpstore.EnvelopeStoreClient
	conn     *grpc.ClientConn
	initOnce = &sync.Once{}
)

func Init(ctx context.Context) {
	initOnce.Do(func() {
		if client != nil {
			return
		}

		var err error
		conn, err = grpcclient.DialContextWithDefaultOptions(
			ctx,
			viper.GetString(grpcTargetEnvelopeStoreViperKey),
		)
		if err != nil {
			e := errors.FromError(err).ExtendComponent(component)
			log.WithError(e).Fatalf("%s: failed to dial grpc server", component)
		}

		client = evlpstore.NewEnvelopeStoreClient(conn)

		log.WithFields(log.Fields{
			"grpc.target": viper.GetString(grpcTargetEnvelopeStoreViperKey),
		}).Infof("%s: client ready", component)
	})
}

func Close() {
	_ = conn.Close()
}

func GlobalEnvelopeStoreClient() evlpstore.EnvelopeStoreClient {
	return client
}

// SetGlobalConfig sets Sarama global configuration
func SetGlobalEnvelopeStoreClient(c evlpstore.EnvelopeStoreClient) {
	client = c
}
