package metrics

import (
	kitmetrics "github.com/go-kit/kit/metrics"
	pkgmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics"
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
