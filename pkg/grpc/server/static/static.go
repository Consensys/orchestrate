package static

import (
	"context"
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/config/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/interceptor"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/server"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/service"
	"google.golang.org/grpc"
)

type Builder struct {
	Options     server.OptionsBuilder
	Interceptor interceptor.Builder
	Service     service.Builder
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Build(ctx context.Context, name string, configuration interface{}) (*grpc.Server, error) {
	cfg, ok := configuration.(*static.Configuration)
	if !ok {
		return nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	// Create Options
	var serverOpts []grpc.ServerOption
	if b.Options != nil {
		opts, err := b.Options.Build(ctx, name, cfg.Options)
		if err != nil {
			return nil, err
		}
		serverOpts = append(serverOpts, opts...)
	}

	// Create interceptors
	var (
		unaryInterceptors  []grpc.UnaryServerInterceptor
		streamInterceptors []grpc.StreamServerInterceptor
		enhancers          []func(srv *grpc.Server)
	)
	if b.Interceptor != nil {
		for _, intercep := range cfg.Interceptors {
			unaryInterceptor, streamInterceptor, enhancer, err := b.Interceptor.Build(ctx, name, intercep)
			if err != nil {
				return nil, err
			}
			if unaryInterceptor != nil {
				unaryInterceptors = append(unaryInterceptors, unaryInterceptor)
			}

			if streamInterceptor != nil {
				streamInterceptors = append(streamInterceptors, streamInterceptor)
			}

			if enhancer != nil {
				enhancers = append(enhancers, enhancer)
			}
		}
	}

	serverOpts = append(serverOpts, grpc.ChainUnaryInterceptor(unaryInterceptors...), grpc.ChainStreamInterceptor(streamInterceptors...))
	enhancer := func(srv *grpc.Server) {
		for i := len(enhancers) - 1; i >= 0; i-- {
			enhancers[i](srv)
		}
	}

	// Create GRPC services
	var register func(srv *grpc.Server)
	if b.Service != nil && cfg.Services != nil {
		registr, err := b.Service.Build(ctx, name, cfg.Services)
		if err != nil {
			return nil, err
		}
		register = registr
	}

	// Create server
	srv := grpc.NewServer(serverOpts...)

	// Register services
	if register != nil {
		register(srv)
	}

	// Enhance server
	if enhancer != nil {
		enhancer(srv)
	}

	return srv, nil
}
