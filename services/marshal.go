package services

import (
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/trace"
)

// Marshaller are responsible to marshal a protobuffer trace into a message
type Marshaller interface {
	Marshal(pb *tracepb.Trace, msg interface{}) error
}

// Unmarshaller are responsible to unmarshal an input message to a protobuf
type Unmarshaller interface {
	Unmarshal(msg interface{}, pb *tracepb.Trace) error
}
