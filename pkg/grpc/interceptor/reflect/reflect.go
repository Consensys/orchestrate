package grpcreflect

import (
	"context"
	"fmt"
	"reflect"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/interceptor"
	"google.golang.org/grpc"
)

type Builder struct {
	builders map[reflect.Type]interceptor.Builder
}

func NewBuilder() *Builder {
	return &Builder{
		builders: make(map[reflect.Type]interceptor.Builder),
	}
}

func (b *Builder) Build(ctx context.Context, name string, configuration interface{}) (grpc.UnaryServerInterceptor, grpc.StreamServerInterceptor, func(srv *grpc.Server), error) {
	builder, ok := b.builders[reflect.TypeOf(configuration)]
	if !ok {
		return nil, nil, nil, fmt.Errorf("no interceptor builder for configuration of type %T (consider adding one)", configuration)
	}

	return builder.Build(ctx, name, configuration)
}

func (b *Builder) AddBuilder(typ reflect.Type, builder interceptor.Builder) {
	b.builders[typ] = builder
}
