package metrics

import (
	kitmetrics "github.com/go-kit/kit/metrics"
)

type metrics struct {
	startedCounter           kitmetrics.Counter
	handledCounter           kitmetrics.Counter
	streamMsgReceivedCounter kitmetrics.Counter
	streamMsgSentCounter     kitmetrics.Counter
	handledDurationHistogram kitmetrics.Histogram
}

func buildMetrics(
	startedCounter, handledCounter, streamMsgReceivedCounter, streamMsgSentCounter kitmetrics.Counter,
	handledDurationHistogram kitmetrics.Histogram,
) *metrics {
	return &metrics{
		startedCounter:           startedCounter,
		handledCounter:           handledCounter,
		streamMsgReceivedCounter: streamMsgReceivedCounter,
		streamMsgSentCounter:     streamMsgSentCounter,
		handledDurationHistogram: handledDurationHistogram,
	}
}

func (r *metrics) StartedCounter() kitmetrics.Counter {
	return r.startedCounter
}

func (r *metrics) HandledCounter() kitmetrics.Counter {
	return r.handledCounter
}

func (r *metrics) StreamMsgReceivedCounter() kitmetrics.Counter {
	return r.streamMsgReceivedCounter
}

func (r *metrics) StreamMsgSentCounter() kitmetrics.Counter {
	return r.streamMsgSentCounter
}

func (r *metrics) HandledDurationHistogram() kitmetrics.Histogram {
	return r.handledDurationHistogram
}
