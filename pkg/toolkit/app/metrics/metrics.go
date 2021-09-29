package metrics

import (
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/config/dynamic"
	"github.com/prometheus/client_golang/prometheus"
)

//go:generate mockgen -source=metrics.go -destination=mock/mock.go -package=mock

const (
	Namespace = "orchestrate"
)

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
