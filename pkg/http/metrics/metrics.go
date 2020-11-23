package metrics

import (
	"github.com/go-kit/kit/metrics/discard"
	"github.com/prometheus/client_golang/prometheus"
	promcol "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/metrics/collector"
	pkgmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/metrics/multi"
)

const (
	Namespace              = "http"
	RequestsTotal          = "requests_total"
	RequestsTLSTotal       = "requests_tls_total"
	RequestsLatencySeconds = "requests_latency_seconds"
	OpenConnections        = "open_connections"
	RetriesTotal           = "retries_total"
	ServerUp               = "server_up"
)

type httpMulti struct {
	*pkgmetrics.DynamicMulti
	*metrics
}

func NewHTTPMetrics(cfg *Config) HTTPMetrics {
	if cfg == nil {
		cfg = NewDefaultConfig()
	}

	multi := pkgmetrics.NewDynamicMulti(nil)

	reqsCounter := promcol.NewCounter(
		multi.Values(),
		&prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      RequestsTotal,
			Help:      "Total count of requests processed on a HTTP service",
		},
		[]string{"tenant_id", "entrypoint", "protocol", "service", "method", "code"},
	)
	multi.Collectors = append(multi.Collectors, reqsCounter.Describe)

	tlsReqsCounter := promcol.NewCounter(
		multi.Values(),
		&prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      RequestsTLSTotal,
			Help:      "Total count of TLS requests processed on a HTTP service",
		},
		[]string{"tenant_id", "entrypoint", "service", "tls_version", "tls_cipher"},
	)
	multi.Collectors = append(multi.Collectors, tlsReqsCounter.Describe)

	reqsLatencyHistogram := promcol.NewHistogram(
		multi.Values(),
		&prometheus.HistogramOpts{
			Namespace: Namespace,
			Name:      RequestsLatencySeconds,
			Help:      "Histogram of service's response latency (seconds)",
			Buckets:   cfg.Buckets,
		},
		[]string{"tenant_id", "entrypoint", "protocol", "service", "method", "code"},
	)
	multi.Collectors = append(multi.Collectors, reqsLatencyHistogram.Describe)

	openConnsGauge := promcol.NewGauge(
		multi.Values(),
		&prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      OpenConnections,
			Help:      "Current count of open connections on a service",
		},
		[]string{"tenant_id", "entrypoint", "protocol", "service", "method"},
	)
	multi.Collectors = append(multi.Collectors, openConnsGauge.Describe)

	retriesCounter := promcol.NewCounter(
		multi.Values(),
		&prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      RetriesTotal,
			Help:      "Total count of request retries on a service.",
		},
		[]string{"tenant_id", "entrypoint", "service"},
	)
	multi.Collectors = append(multi.Collectors, retriesCounter.Describe)

	serverUpGauge := promcol.NewGauge(
		multi.Values(),
		&prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      ServerUp,
			Help:      "Current server status (0=DOWN, 1=UP)",
		},
		[]string{"tenant_id", "entrypoint", "service", "url"},
	)
	multi.Collectors = append(multi.Collectors, serverUpGauge.Describe)

	return &httpMulti{
		DynamicMulti: multi,
		metrics: buildMetrics(
			reqsCounter,
			tlsReqsCounter,
			reqsLatencyHistogram,
			openConnsGauge,
			retriesCounter,
			serverUpGauge,
		),
	}
}

func NewHTTPNopMetrics(_ *Config) HTTPMetrics {
	return &httpMulti{
		DynamicMulti: pkgmetrics.NewDynamicMulti(nil),
		metrics: buildMetrics(
			discard.NewCounter(),
			discard.NewCounter(),
			discard.NewHistogram(),
			discard.NewGauge(),
			discard.NewCounter(),
			discard.NewGauge(),
		),
	}
}
