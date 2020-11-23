package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/config/dynamic"
)

//go:generate mockgen -source=metrics.go -destination=mock/mock.go -package=mock

type Prometheus interface {
	Describe(chan<- *prometheus.Desc)
	Collect(chan<- prometheus.Metric)
}

type DynamicPrometheus interface {
	Switch(*dynamic.Configuration) error
	Prometheus
}

type Registry interface {
	Add(m Prometheus)
	SwitchDynConfig(dynCfg *dynamic.Configuration) error
	Prometheus
}
