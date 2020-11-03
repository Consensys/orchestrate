package collector

import (
	kitmetrics "github.com/go-kit/kit/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type Histogram struct {
	*Collector
}

func NewHistogram(values chan<- *Collector, opts *prometheus.HistogramOpts, labelNames []string) *Histogram {
	return &Histogram{
		NewCollector(
			"", opts.Name,
			nil,
			prometheus.NewHistogramVec(*opts, labelNames),
			nil,
			values,
		),
	}
}

func (h *Histogram) With(labelValues ...string) kitmetrics.Histogram {
	return &Histogram{h.Collector.withLabelValues(labelValues...)}
}
