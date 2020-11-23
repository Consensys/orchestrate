package dynamic

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/handler"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/handler/healthcheck"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/handler/prometheus"
	reflecthandler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/handler/reflect"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/handler/swagger"
)

type Builder struct {
	reflect *reflecthandler.Builder
}

func NewBuilder() *Builder {
	b := &Builder{
		reflect: reflecthandler.NewBuilder(),
	}

	b.AddBuilder(reflect.TypeOf(&dynamic.HealthCheck{}), healthcheck.NewTraefikBuilder())
	b.AddBuilder(reflect.TypeOf(&dynamic.Swagger{}), swagger.NewBuilder())
	b.AddBuilder(reflect.TypeOf(&dynamic.Prometheus{}), prometheus.NewBuilder(nil))

	return b
}

func (b *Builder) AddBuilder(typ reflect.Type, builder handler.Builder) {
	b.reflect.AddBuilder(typ, builder)
}

func (b *Builder) Build(ctx context.Context, name string, configuration interface{}, respModifier func(*http.Response) error) (http.Handler, error) {
	cfg, ok := configuration.(*dynamic.Service)
	if !ok {
		return nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	field, err := cfg.Field()
	if err != nil {
		return nil, err
	}

	return b.reflect.Build(ctx, name, field, respModifier)
}
