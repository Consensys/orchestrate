package standard

import (
	kitmetrics "github.com/go-kit/kit/metrics"
)

type GRPCServer struct {
	startedCounter           kitmetrics.Counter
	handledCounter           kitmetrics.Counter
	streamMsgReceivedCounter kitmetrics.Counter
	streamMsgSentCounter     kitmetrics.Counter
	handledDurationHistogram kitmetrics.Histogram
}

func NewGRPCServer(
	startedCounter, handledCounter, streamMsgReceivedCounter, streamMsgSentCounter kitmetrics.Counter,
	handledDurationHistogram kitmetrics.Histogram,
) *GRPCServer {
	return &GRPCServer{
		startedCounter:           startedCounter,
		handledCounter:           handledCounter,
		streamMsgReceivedCounter: streamMsgReceivedCounter,
		streamMsgSentCounter:     streamMsgSentCounter,
		handledDurationHistogram: handledDurationHistogram,
	}
}

func (r *GRPCServer) StartedCounter() kitmetrics.Counter {
	return r.startedCounter
}

func (r *GRPCServer) HandledCounter() kitmetrics.Counter {
	return r.handledCounter
}

func (r *GRPCServer) StreamMsgReceivedCounter() kitmetrics.Counter {
	return r.streamMsgReceivedCounter
}

func (r *GRPCServer) StreamMsgSentCounter() kitmetrics.Counter {
	return r.streamMsgSentCounter
}

func (r *GRPCServer) HandledDurationHistogram() kitmetrics.Histogram {
	return r.handledDurationHistogram
}
