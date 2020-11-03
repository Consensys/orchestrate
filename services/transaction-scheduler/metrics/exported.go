package metrics

import (
	kitmetrics "github.com/go-kit/kit/metrics"
	pkgmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics"
)

//go:generate mockgen -source=exported.go -destination=mock/mock.go -package=mock

const ModuleName = "transaction_scheduler"

type TransactionSchedulerMetrics interface {
	CreatedJobsCounter() kitmetrics.Counter
	pkgmetrics.Prometheus
}
