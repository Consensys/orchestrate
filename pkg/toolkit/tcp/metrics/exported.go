package metrics

import (
	pkgmetrics "github.com/ConsenSys/orchestrate/pkg/toolkit/app/metrics"
	kitmetrics "github.com/go-kit/kit/metrics"
)

//go:generate mockgen -source=exported.go -destination=mock/mock.go -package=mock

const ModuleName = "tcp"

type TPCMetrics interface {
	AcceptedConnsCounter() kitmetrics.Counter
	ClosedConnsCounter() kitmetrics.Counter
	ConnsLatencyHistogram() kitmetrics.Histogram
	OpenConnsGauge() kitmetrics.Gauge
	pkgmetrics.Prometheus
}
