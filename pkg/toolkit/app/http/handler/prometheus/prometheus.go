package prometheus

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Builder struct {
	registry prometheus.Gatherer
}

func NewBuilder(registry prometheus.Gatherer) *Builder {
	if registry == nil {
		registry = prometheus.DefaultGatherer
	}

	return &Builder{
		registry: registry,
	}
}

func (b *Builder) Build(ctx context.Context, name string, configuration interface{}, respModifier func(resp *http.Response) error) (http.Handler, error) {
	return promhttp.HandlerFor(b.registry, promhttp.HandlerOpts{}), nil
}
