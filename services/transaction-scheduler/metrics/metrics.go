package metrics

import (
	"github.com/go-kit/kit/metrics/discard"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	pkgmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/metrics/multi"
)

const (
	Namespace      = "transaction_scheduler"
	CreatedJobName = "created_job"
)

type tpcMetrics struct {
	prometheus.Collector
	*metrics
}

func NewTransactionSchedulerMetrics() TransactionSchedulerMetrics {
	multi := pkgmetrics.NewMulti()

	createdJobsCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      CreatedJobName,
			Help:      "Total count of created jobs.",
		},
		[]string{"tenant_id", "chain_uuid"},
	)
	multi.Collectors = append(multi.Collectors, createdJobsCounter)

	return &tpcMetrics{
		Collector: multi,
		metrics: buildMetrics(
			kitprometheus.NewCounter(createdJobsCounter),
		),
	}
}

func NewTransactionSchedulerNopMetrics() TransactionSchedulerMetrics {
	return &tpcMetrics{
		Collector: pkgmetrics.NewMulti(),
		metrics: buildMetrics(
			discard.NewCounter(),
		),
	}
}
