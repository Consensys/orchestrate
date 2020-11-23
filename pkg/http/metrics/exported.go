package metrics

import (
	kitmetrics "github.com/go-kit/kit/metrics"
	pkgmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/metrics"
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
