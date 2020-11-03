package collector

import (
	kitmetrics "github.com/go-kit/kit/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type Gauge struct {
	*Collector
}

func NewGauge(values chan<- *Collector, opts *prometheus.GaugeOpts, labelNames []string) *Gauge {
	g := &Gauge{
		NewCollector(
			"", opts.Name,
			nil,
			prometheus.NewGaugeVec(*opts, labelNames),
			nil,
			values,
		),
	}

	if len(labelNames) == 0 {
		g.Set(0)
	}

	return g
}

func (g *Gauge) With(labelValues ...string) kitmetrics.Gauge {
	return &Gauge{g.Collector.withLabelValues(labelValues...)}
}
