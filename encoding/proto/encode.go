package proto

import (
	"github.com/golang/protobuf/proto"
	errors "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
)

// Marshal a proto into a message assumed to be an Envelope
func Marshal(pb proto.Message) ([]byte, error) {
	buf, err := proto.Marshal(pb)
	if err != nil {
		return nil, errors.EncodingError(err.Error()).SetComponent(component)
	}
	return buf, nil
}
