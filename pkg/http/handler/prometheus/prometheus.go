package prometheus

import (
	"context"
	"net/http"

	"github.com/containous/traefik/v2/pkg/metrics"
	traefiktypes "github.com/containous/traefik/v2/pkg/types"
)

type Builder struct{}

func NewBuilder(cfg *traefiktypes.Prometheus) *Builder {
	metrics.RegisterPrometheus(context.Background(), cfg)
	return &Builder{}
}

func (b *Builder) Build(ctx context.Context, name string, configuration interface{}, respModifier func(resp *http.Response) error) (http.Handler, error) {
	return metrics.PrometheusHandler(), nil
}
