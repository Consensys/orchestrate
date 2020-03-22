package client

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/opentracing"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/tracing/opentracing/jaeger"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
	gogrpc "google.golang.org/grpc"
)

const component = "envelope-store.client"

type serviceName string

var (
	client   svc.EnvelopeStoreClient
	conn     *gogrpc.ClientConn
	initOnce = &sync.Once{}
)

func Init(ctx context.Context) {
	initOnce.Do(func() {
		if client != nil {
			return
		}

		ctxWithValue := context.WithValue(ctx, serviceName("service-name"), viper.GetString(jaeger.ServiceNameViperKey))
		opentracing.Init(ctxWithValue)
		var err error
		conn, err = grpc.DialContextWithDefaultOptions(
			ctx,
			viper.GetString(EnvelopeStoreURLViperKey),
		)
		if err != nil {
			e := errors.FromError(err).ExtendComponent(component)
			log.WithError(e).Fatalf("%s: failed to dial grpc server", component)
		}

		client = svc.NewEnvelopeStoreClient(conn)

		log.WithFields(log.Fields{
			"url": viper.GetString(EnvelopeStoreURLViperKey),
		}).Infof("%s: client ready", component)
	})
}

func Close() {
	_ = conn.Close()
}

func GlobalEnvelopeStoreClient() svc.EnvelopeStoreClient {
	return client
}

// SetGlobalConfig sets Sarama global configuration
func SetGlobalEnvelopeStoreClient(c svc.EnvelopeStoreClient) {
	client = c
}
