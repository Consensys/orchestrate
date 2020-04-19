package prometheus

import (
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/standard"
)

const tcpMetricsNamespace = "tcp"

type TCP struct {
	prometheus.Collector
	*standard.TCP
}

func NewTCP(cfg *Config) *TCP {
	multi := NewMulti()

	acceptedConnsCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: tcpMetricsNamespace,
			Name:      "accepted_conns_total",
			Help:      "Total count of accepted connections.",
		},
		[]string{"entrypoint"},
	)
	multi.collectors = append(multi.collectors, acceptedConnsCounter)

	closedConnsCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: tcpMetricsNamespace,
			Name:      "closed_conns_total",
			Help:      "Total count of closed connections.",
		},
		[]string{"entrypoint"},
	)
	multi.collectors = append(multi.collectors, closedConnsCounter)

	connsLatencyHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: tcpMetricsNamespace,
			Name:      "open_conns_duration_seconds",
			Help:      "Histogram of connections duration (seconds)",
			Buckets:   cfg.TCP.Buckets,
		},
		[]string{"entrypoint"},
	)
	multi.collectors = append(multi.collectors, connsLatencyHistogram)

	openConnsGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: tcpMetricsNamespace,
			Name:      "open_conns",
			Help:      "Current count of open connections on an entrypoint",
		},
		[]string{"entrypoint"},
	)
	multi.collectors = append(multi.collectors, openConnsGauge)

	return &TCP{
		Collector: multi,
		TCP: standard.NewTCP(
			kitprometheus.NewCounter(acceptedConnsCounter),
			kitprometheus.NewCounter(closedConnsCounter),
			kitprometheus.NewHistogram(connsLatencyHistogram),
			kitprometheus.NewGauge(openConnsGauge),
		),
	}
}
