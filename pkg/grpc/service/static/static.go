package static

import (
	"context"
	"fmt"
	"reflect"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/grpc/config/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/grpc/service"
	grpcreflect "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/grpc/service/reflect"
	"google.golang.org/grpc"
)

type Builder struct {
	reflect *grpcreflect.Builder
}

func NewBuilder() *Builder {
	builder := grpcreflect.NewBuilder()
	return &Builder{
		reflect: builder,
	}
}

func (b *Builder) AddBuilder(typ reflect.Type, builder service.Builder) {
	b.reflect.AddBuilder(typ, builder)
}

func (b *Builder) Build(ctx context.Context, name string, configuration interface{}) (func(srv *grpc.Server), error) {
	cfg, ok := configuration.(*static.Services)
	if !ok {
		return nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	field, err := cfg.Field()
	if err != nil {
		return nil, err
	}

	return b.reflect.Build(ctx, name, field)
}
