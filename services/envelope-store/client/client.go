package client

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/opentracing"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
	gogrpc "google.golang.org/grpc"
)

const component = "envelope-store.client"

type serviceName string
type Client struct {
	srv  svc.EnvelopeStoreClient
	conn *gogrpc.ClientConn
}

func NewClient(ctx context.Context, cfg *Config) (*Client, error) {
	ctxWithValue := context.WithValue(ctx, serviceName("service-name"), cfg.serviceName)
	opentracing.Init(ctxWithValue)

	conn, err := grpc.DialContextWithDefaultOptions(
		ctx,
		cfg.envelopeStoreURL,
	)
	if err != nil {
		e := errors.FromError(err).ExtendComponent(component)
		log.WithError(e).Fatalf("%s: failed to dial grpc server", component)
		return nil, err
	}

	client := svc.NewEnvelopeStoreClient(conn)

	log.WithFields(log.Fields{
		"url": cfg.envelopeStoreURL,
	}).Infof("%s: client ready", component)

	return &Client{
		srv:  client,
		conn: conn,
	}, nil
}
