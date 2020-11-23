package reflect

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/middleware"
)

type Builder struct {
	builders map[reflect.Type]middleware.Builder
}

func NewBuilder() *Builder {
	return &Builder{
		builders: make(map[reflect.Type]middleware.Builder),
	}
}

func (b *Builder) Build(ctx context.Context, name string, configuration interface{}) (mid func(http.Handler) http.Handler, respModifier func(resp *http.Response) error, err error) {
	log.FromContext(ctx).
		WithField("middleware", name).
		WithField("type", fmt.Sprintf("%T", configuration)).
		Debugf("building middleware")

	builder, ok := b.builders[reflect.TypeOf(configuration)]
	if !ok {
		return nil, nil, fmt.Errorf("no middleware builder for configuration of type %T (consider adding one)", configuration)
	}

	return builder.Build(ctx, name, configuration)
}

func (b *Builder) AddBuilder(typ reflect.Type, builder middleware.Builder) {
	b.builders[typ] = builder
}
