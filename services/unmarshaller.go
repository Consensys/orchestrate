package services

import (
	proto "github.com/golang/protobuf/proto"
)

// Unmarshaller are responsible to unmarshal high level input message into a protobuf message
type Unmarshaller interface {
	// Unmarshal high level input message into a protobuf message
	Unmarshal(msg interface{}, pb proto.Message) error
}
