package contractregistry

import (
	"reflect"

	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/config/static"
	grpcauth "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/interceptor/auth"
	grpclogrus "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/interceptor/logrus"
	grpcmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/interceptor/metrics"
	staticinterceptor "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/interceptor/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/server"
	staticserver "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/server/static"
	staticservice "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/service/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/proto"
	grpcservice "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/service/grpc"
)

func NewGRPCBuilder(
	service svc.ContractRegistryServer,
	checker auth.Checker, multitenancy bool,
	logger *logrus.Logger,
	reg metrics.GRPCServer,
) (server.Builder, error) {
	// Create GRPC server builder
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

	// Add Builder for Metrics interceptor
	interceptorBuilder.AddBuilder(
		reflect.TypeOf(&static.Metrics{}),
		grpcmetrics.NewBuilder(reg),
	)

	builder.Interceptor = interceptorBuilder

	// Add Builder for Contract-Registry service
	serviceBuilder := staticservice.NewBuilder()
	serviceBuilder.AddBuilder(
		reflect.TypeOf(&static.Contracts{}),
		grpcservice.NewBuilder(service),
	)
	builder.Service = serviceBuilder

	return builder, nil
}

func NewGRPCStaticConfig() *static.Configuration {
	return &static.Configuration{
		Services: &static.Services{
			Contracts: &static.Contracts{},
		},
		Interceptors: []*static.Interceptor{
			{Tags: &static.Tags{}},
			{Logrus: &static.Logrus{}},
			{Auth: &static.Auth{}},
			{Error: &static.Error{}},
			{Recovery: &static.Recovery{}},
		},
	}
}
