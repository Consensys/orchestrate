package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Multi struct {
	collectors []prometheus.Collector
}

func NewMulti() *Multi {
	return &Multi{}
}

// Describe implements prometheus.Collector and simply calls
// the registered describer functions.
func (c *Multi) Describe(ch chan<- *prometheus.Desc) {
	for _, cllctr := range c.collectors {
		cllctr.Describe(ch)
	}
}

// Collect collectors
func (c *Multi) Collect(ch chan<- prometheus.Metric) {
	for _, cllctr := range c.collectors {
		cllctr.Collect(ch)
	}
}
