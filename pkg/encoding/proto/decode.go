package proto

import (
	"github.com/ConsenSys/orchestrate/pkg/errors"
	"google.golang.org/protobuf/proto"
)

// Unmarshal parses the protocol buffer representation in `buf`
// and places the decoded result in `pb`
//
// Unmarshal resets pb before starting to unmarshal
func Unmarshal(buf []byte, pb proto.Message) error {
	// Unmarshal
	err := proto.Unmarshal(buf, pb)
	if err != nil {
		return errors.EncodingError(err.Error()).SetComponent(component)
	}

	return nil
}

// UnmarshalMerge parses the protocol buffer representation in buf and
// writes the decoded result to pb
func UnmarshalMerge(buf []byte, pb proto.Message) error {
	// Unmarshal
	err := proto.UnmarshalOptions{Merge: true}.Unmarshal(buf, pb)
	if err != nil {
		return errors.EncodingError(err.Error()).SetComponent(component)
	}

	return nil
}
