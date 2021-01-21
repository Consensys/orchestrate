package reflect

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/handler"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"
)

const component = "http.handler"

type Builder struct {
	builders map[reflect.Type]handler.Builder
}

func NewBuilder() *Builder {
	return &Builder{
		builders: make(map[reflect.Type]handler.Builder),
	}
}

func (b *Builder) Build(ctx context.Context, name string, configuration interface{}, respModifier func(resp *http.Response) error) (http.Handler, error) {
	log.NewLogger().WithContext(ctx).SetComponent(component).
		WithField("handler_name", name).
		WithField("type", fmt.Sprintf("%T", configuration)).
		Debug("building handler")

	builder, ok := b.builders[reflect.TypeOf(configuration)]
	if !ok {
		return nil, fmt.Errorf("no service builder for configuration of type %T (consider adding one)", configuration)
	}

	return builder.Build(ctx, name, configuration, respModifier)
}

func (b *Builder) AddBuilder(typ reflect.Type, builder handler.Builder) {
	b.builders[typ] = builder
}
