package services

import (
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/trace"
)

// OffsetMarker is an interface to mark that a context has been processed
// Typically marking kafka offsets of message that have been consumed
type OffsetMarker interface {
	Mark(msg interface{}) error
}

// TraceProducer produces a trace in another service typically a Kafka queue
type TraceProducer interface {
	Produce(pb *tracepb.Trace) error
}
