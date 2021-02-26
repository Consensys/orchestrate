package metrics

import (
	pkgmetrics "github.com/ConsenSys/orchestrate/pkg/toolkit/app/metrics"
	kitmetrics "github.com/go-kit/kit/metrics"
)

//go:generate mockgen -source=exported.go -destination=mock/mock.go -package=mock

const ModuleName = "http"

type HTTPMetrics interface {
	RequestsCounter() kitmetrics.Counter
	TLSRequestsCounter() kitmetrics.Counter
	RequestsLatencyHistogram() kitmetrics.Histogram
	OpenConnsGauge() kitmetrics.Gauge
	RetriesCounter() kitmetrics.Counter
	ServerUpGauge() kitmetrics.Gauge
	pkgmetrics.DynamicPrometheus
}
