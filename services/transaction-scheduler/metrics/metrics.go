package metrics

import (
	"github.com/go-kit/kit/metrics/discard"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	metrics1 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/metrics"
	pkgmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/metrics/multi"
)

const (
	Subsystem           = "transaction_scheduler"
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
