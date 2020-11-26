package metrics

import (
	kitmetrics "github.com/go-kit/kit/metrics"
)

type metrics struct {
	blockCounter kitmetrics.Counter
}

// IsSync metric

func buildMetrics(
	blockCounter kitmetrics.Counter,
) *metrics {
	return &metrics{
		blockCounter: blockCounter,
	}
}

func (r *metrics) BlockCounter() kitmetrics.Counter {
	return r.blockCounter
}
