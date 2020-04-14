package grpc

import (
	"context"
	"reflect"

	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/config/static"
	"google.golang.org/grpc"

	grpcauth "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/interceptor/auth"
	grpclogrus "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/interceptor/logrus"
	staticinterceptor "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/interceptor/static"
	staticserver "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/server/static"
	staticservice "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/service/static"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
)

type ServerBuilder struct {
	*staticserver.Builder
}

func NewServerBuilder(
	srv svc.EnvelopeStoreServer,
	checker auth.Checker,
	multitenancy bool,
	logger *logrus.Logger,
) (ServerBuilder, error) {
	builder := staticserver.NewBuilder()

	// Create interceptor builder
	interceptorBuilder := staticinterceptor.NewBuilder()

	// Add Builder for Authentication interceptor
	interceptorBuilder.AddBuilder(
		reflect.TypeOf(&static.Auth{}),
		grpcauth.NewBuilder(checker, multitenancy),
	)

	// Add Builder for Authentication interceptor
	interceptorBuilder.AddBuilder(
		reflect.TypeOf(&static.Logrus{}),
		grpclogrus.NewBuilder(logger, logrus.Fields{"system": "grpc.internal"}),
	)

	builder.Interceptor = interceptorBuilder

	// Add Builder for Envelope-store service
	serviceBuilder := staticservice.NewBuilder()
	serviceBuilder.AddBuilder(
		reflect.TypeOf(&static.Envelopes{}),
		newServiceBuilder(srv),
	)
	builder.Service = serviceBuilder

	return ServerBuilder{
		builder,
	}, nil
}

func (b *ServerBuilder) BuildServer(ctx context.Context, name string, cfg interface{}) (*grpc.Server, error) {
	return b.Build(ctx, name, cfg)
}
