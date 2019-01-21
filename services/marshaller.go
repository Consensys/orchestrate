package services

import (
	proto "github.com/golang/protobuf/proto"
)

// Marshaller are responsible to marshal a protobuffer message to a higher level message format
type Marshaller interface {
	// Marshal a protobuffer message to a higher level message format
	Marshal(pb proto.Message, msg interface{}) error
}
