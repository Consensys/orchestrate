package proto

import (
	"github.com/consensys/orchestrate/pkg/errors"
	"google.golang.org/protobuf/proto"
)

// Marshal a proto into a message assumed to be an Envelope
func Marshal(pb proto.Message) ([]byte, error) {
	buf, err := proto.Marshal(pb)
	if err != nil {
		return nil, errors.EncodingError(err.Error()).SetComponent(component)
	}
	return buf, nil
}
