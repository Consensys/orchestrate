package metrics

import (
	kitmetrics "github.com/go-kit/kit/metrics"
)

type metrics struct {
	createdJobsCounter kitmetrics.Counter
}

func buildMetrics(
	createdJobsCounter kitmetrics.Counter,
) *metrics {
	return &metrics{
		createdJobsCounter: createdJobsCounter,
	}
}

func (r *metrics) CreatedJobsCounter() kitmetrics.Counter {
	return r.createdJobsCounter
}
