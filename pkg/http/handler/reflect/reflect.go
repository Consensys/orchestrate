package reflect

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/handler"
)

type Builder struct {
	builders map[reflect.Type]handler.Builder
}

func NewBuilder() *Builder {
	return &Builder{
		builders: make(map[reflect.Type]handler.Builder),
	}
}

func (b *Builder) Build(ctx context.Context, name string, configuration interface{}, respModifier func(resp *http.Response) error) (http.Handler, error) {
	log.FromContext(ctx).
		WithField("handler", name).
		WithField("type", fmt.Sprintf("%T", configuration)).
		Debugf("building handler")

	builder, ok := b.builders[reflect.TypeOf(configuration)]
	if !ok {
		return nil, fmt.Errorf("no service builder for configuration of type %T (consider adding one)", configuration)
	}

	return builder.Build(ctx, name, configuration, respModifier)
}

func (b *Builder) AddBuilder(typ reflect.Type, builder handler.Builder) {
	b.builders[typ] = builder
}
