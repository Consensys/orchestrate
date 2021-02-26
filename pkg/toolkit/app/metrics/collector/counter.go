package collector

import (
	kitmetrics "github.com/go-kit/kit/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type Counter struct {
	*Collector
}

func NewCounter(values chan<- *Collector, opts *prometheus.CounterOpts, labelNames []string) *Counter {
	c := &Counter{
		NewCollector(
			"", opts.Name,
			nil,
			prometheus.NewCounterVec(*opts, labelNames),
			nil,
			values,
		),
	}

	if len(labelNames) == 0 {
		c.Add(0)
	}

	return c
}

func (c *Counter) With(labelValues ...string) kitmetrics.Counter {
	return &Counter{c.Collector.withLabelValues(labelValues...)}
}
