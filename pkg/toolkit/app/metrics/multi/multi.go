package multi

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Multi struct {
	Collectors []prometheus.Collector
}

func NewMulti() *Multi {
	return &Multi{}
}

// Describe implements prometheus.Collector and simply calls
// the registered describer functions.
func (c *Multi) Describe(ch chan<- *prometheus.Desc) {
	for _, cllctr := range c.Collectors {
		cllctr.Describe(ch)
	}
}

// Collect collectors
func (c *Multi) Collect(ch chan<- prometheus.Metric) {
	for _, cllctr := range c.Collectors {
		cllctr.Collect(ch)
	}
}
