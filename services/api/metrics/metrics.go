package metrics

import (
	metrics1 "github.com/ConsenSys/orchestrate/pkg/metrics"
	pkgmetrics "github.com/ConsenSys/orchestrate/pkg/metrics/multi"
	"github.com/go-kit/kit/metrics/discard"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	Subsystem           = "api"
	JobLatencySeconds   = "job_latency_seconds"
	MinedLatencySeconds = "mined_latency_seconds"
)

type tpcMetrics struct {
	prometheus.Collector
	*metrics
}

func NewTransactionSchedulerMetrics() TransactionSchedulerMetrics {
	multi := pkgmetrics.NewMulti()

	jobsLatencyHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: metrics1.Namespace,
			Subsystem: Subsystem,
			Name:      JobLatencySeconds,
			Help:      "Histogram of job latency between status (second). Except PENDING and MINED, see mined_latency_seconds",
			Buckets:   []float64{.01, .025, .05, .1, 1, 5},
		},
		[]string{"chain_uuid", "status", "prev_status"},
	)
	multi.Collectors = append(multi.Collectors, jobsLatencyHistogram)

	minedLatencyHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: metrics1.Namespace,
			Subsystem: Subsystem,
			Name:      MinedLatencySeconds,
			Help:      "Histogram of latency between PENDING and MINED (second)",
			Buckets:   []float64{.5, 1, 5, 10, 15, 20},
		},
		[]string{"chain_uuid", "status", "prev_status"},
	)
	multi.Collectors = append(multi.Collectors, minedLatencyHistogram)

	return &tpcMetrics{
		Collector: multi,
		metrics: buildMetrics(
			kitprometheus.NewHistogram(jobsLatencyHistogram),
			kitprometheus.NewHistogram(minedLatencyHistogram),
		),
	}
}

func NewTransactionSchedulerNopMetrics() TransactionSchedulerMetrics {
	return &tpcMetrics{
		Collector: pkgmetrics.NewMulti(),
		metrics: buildMetrics(
			discard.NewHistogram(),
			discard.NewHistogram(),
		),
	}
}
