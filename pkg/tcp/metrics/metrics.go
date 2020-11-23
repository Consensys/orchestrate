package metrics

import (
	"github.com/go-kit/kit/metrics/discard"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	pkgmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/metrics/multi"
)

const (
	Namespace                = "tcp"
	AcceptedConnsTotal       = "accepted_conns_total"
	ClosedConnsTotal         = "closed_conns_total"
	OpenConnsDurationSeconds = "open_conns_duration_seconds"
	OpenConns                = "open_conns"
)

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
			Namespace: Namespace,
			Name:      AcceptedConnsTotal,
			Help:      "Total count of accepted connections.",
		},
		[]string{"entrypoint"},
	)
	multi.Collectors = append(multi.Collectors, acceptedConnsCounter)

	closedConnsCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      ClosedConnsTotal,
			Help:      "Total count of closed connections.",
		},
		[]string{"entrypoint"},
	)
	multi.Collectors = append(multi.Collectors, closedConnsCounter)

	connsLatencyHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: Namespace,
			Name:      OpenConnsDurationSeconds,
			Help:      "Histogram of connections duration (seconds)",
			Buckets:   cfg.Buckets,
		},
		[]string{"entrypoint"},
	)
	multi.Collectors = append(multi.Collectors, connsLatencyHistogram)

	openConnsGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      OpenConns,
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
