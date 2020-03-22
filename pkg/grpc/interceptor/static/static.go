package static

import (
	"context"
	"fmt"
	"reflect"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/config/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/interceptor"
	grpcerror "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/interceptor/error"
	grpcrecovery "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/interceptor/recovery"
	grpcreflect "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/interceptor/reflect"
	grpctags "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/interceptor/tags"
	"google.golang.org/grpc"
)

type Builder struct {
	reflect *grpcreflect.Builder
}

func NewBuilder() *Builder {
	builder := grpcreflect.NewBuilder()
	builder.AddBuilder(reflect.TypeOf(&static.Error{}), grpcerror.NewBuilder())
	builder.AddBuilder(reflect.TypeOf(&static.Recovery{}), grpcrecovery.NewBuilder())
	builder.AddBuilder(reflect.TypeOf(&static.Tags{}), grpctags.NewBuilder())

	return &Builder{
		reflect: builder,
	}
}

func (b *Builder) AddBuilder(typ reflect.Type, builder interceptor.Builder) {
	b.reflect.AddBuilder(typ, builder)
}

func (b *Builder) Build(ctx context.Context, name string, configuration interface{}) (grpc.UnaryServerInterceptor, grpc.StreamServerInterceptor, func(srv *grpc.Server), error) {
	cfg, ok := configuration.(*static.Interceptor)
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid interceptor configuration type (expected %T but got %T)", cfg, configuration)
	}

	field, err := cfg.Field()
	if err != nil {
		return nil, nil, nil, err
	}

	return b.reflect.Build(ctx, name, field)
}
