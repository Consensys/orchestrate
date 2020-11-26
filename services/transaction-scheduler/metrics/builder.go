package metrics

import (
	kitmetrics "github.com/go-kit/kit/metrics"
)

type metrics struct {
	jobsLatencyHistogram  kitmetrics.Histogram
	minedLatencyHistogram kitmetrics.Histogram
}

func buildMetrics(
	jobsLatencyHistogram,
	minedLatencyHistogram kitmetrics.Histogram,
) *metrics {
	return &metrics{
		jobsLatencyHistogram:  jobsLatencyHistogram,
		minedLatencyHistogram: minedLatencyHistogram,
	}
}

func (r *metrics) JobsLatencyHistogram() kitmetrics.Histogram {
	return r.jobsLatencyHistogram
}

func (r *metrics) MinedLatencyHistogram() kitmetrics.Histogram {
	return r.minedLatencyHistogram
}
