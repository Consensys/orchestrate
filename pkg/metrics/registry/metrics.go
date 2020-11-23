package registry

import (
	prom "github.com/prometheus/client_golang/prometheus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/metrics"
)

type metricRegistry struct {
	metrics.Prometheus
	metrics []metrics.Prometheus
}

func NewMetricRegistry(mecs ...metrics.Prometheus) metrics.Registry {
	return &metricRegistry{
		metrics: mecs,
	}
}

func (reg *metricRegistry) Add(m metrics.Prometheus) {
	reg.metrics = append(reg.metrics, m)
}

func (reg *metricRegistry) SwitchDynConfig(dynCfg *dynamic.Configuration) error {
	for _, m := range reg.metrics {
		if http, ok := m.(metrics.DynamicPrometheus); ok {
			if err := http.Switch(dynCfg); err != nil {
				return err
			}
		}
	}

	return nil
}

// Describe implements prometheus.Collector and simply calls
// the registered describer functions.
func (reg *metricRegistry) Describe(ch chan<- *prom.Desc) {
	for _, m := range reg.metrics {
		m.Describe(ch)
	}
}

// Collect collectors
func (reg *metricRegistry) Collect(ch chan<- prom.Metric) {
	for _, m := range reg.metrics {
		m.Collect(ch)
	}
}
