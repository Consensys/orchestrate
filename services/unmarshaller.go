package services

import (
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core.git/protobuf/trace"
)

// Unmarshaller are responsible to unmarshal an input message to a protobuf
type Unmarshaller interface {
	Unmarshal(msg interface{}, pb *tracepb.Trace) error
}
