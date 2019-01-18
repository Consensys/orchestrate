package services

import (
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core.git/protobuf/trace"
)

// Marshaller are responsible to marshal a protobuffer trace into a message
type Marshaller interface {
	Marshal(pb *tracepb.Trace, msg interface{}) error
}
