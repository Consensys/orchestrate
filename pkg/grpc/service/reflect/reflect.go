package reflect

import (
	"context"
	"fmt"
	"reflect"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/service"
	"google.golang.org/grpc"
)

type Builder struct {
	builders map[reflect.Type]service.Builder
}

func NewBuilder() *Builder {
	return &Builder{
		builders: make(map[reflect.Type]service.Builder),
	}
}

func (b *Builder) AddBuilder(typ reflect.Type, builder service.Builder) {
	b.builders[typ] = builder
}

func (b *Builder) Build(ctx context.Context, name string, configuration interface{}) (func(srv *grpc.Server), error) {
	builder, ok := b.builders[reflect.TypeOf(configuration)]
	if !ok {
		return nil, fmt.Errorf("no service builder for configuration of type %T (consider adding one)", configuration)
	}

	return builder.Build(ctx, name, configuration)
}
