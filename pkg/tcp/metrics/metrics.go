package metrics

import (
	"github.com/go-kit/kit/metrics/discard"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	pkgmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/multi"
)

const metricsNamespace = "tcp"

type tpcMetrics struct {
	prometheus.Collector
	*metrics
}

func NewTCPMetrics(cfg *Config) TPCMetrics {
	if cfg == nil {
		cfg = NewDefaultConfig()
	}

	multi := pkgmetrics.NewMulti()

	acceptedConnsCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: metricsNamespace,
			Name:      "accepted_conns_total",
			Help:      "Total count of accepted connections.",
		},
		[]string{"entrypoint"},
	)
	multi.Collectors = append(multi.Collectors, acceptedConnsCounter)

	closedConnsCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: metricsNamespace,
			Name:      "closed_conns_total",
			Help:      "Total count of closed connections.",
		},
		[]string{"entrypoint"},
	)
	multi.Collectors = append(multi.Collectors, closedConnsCounter)

	connsLatencyHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: metricsNamespace,
			Name:      "open_conns_duration_seconds",
			Help:      "Histogram of connections duration (seconds)",
			Buckets:   cfg.Buckets,
		},
		[]string{"entrypoint"},
	)
	multi.Collectors = append(multi.Collectors, connsLatencyHistogram)

	openConnsGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: metricsNamespace,
			Name:      "open_conns",
			Help:      "Current count of open connections on an entrypoint",
		},
		[]string{"entrypoint"},
	)
	multi.Collectors = append(multi.Collectors, openConnsGauge)

	return &tpcMetrics{
		Collector: multi,
		metrics: buildMetrics(
			kitprometheus.NewCounter(acceptedConnsCounter),
			kitprometheus.NewCounter(closedConnsCounter),
			kitprometheus.NewHistogram(connsLatencyHistogram),
			kitprometheus.NewGauge(openConnsGauge),
		),
	}
}

func NewTCPNopMetrics(_ *Config) TPCMetrics {
	return &tpcMetrics{
		Collector: pkgmetrics.NewMulti(),
		metrics: buildMetrics(
			discard.NewCounter(),
			discard.NewCounter(),
			discard.NewHistogram(),
			discard.NewGauge(),
		),
	}
}
