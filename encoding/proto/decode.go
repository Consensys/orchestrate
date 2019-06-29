package proto

import (
	"github.com/golang/protobuf/proto"
	errors "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
)

// Unmarshal parses the protocol buffer representation in `buf` and places the decoded result in `pb`
func Unmarshal(buf []byte, pb proto.Message) error {
	// Cast message into protobuffer
	err := proto.Unmarshal(buf, pb)
	if err != nil {
		return errors.EncodingError(err).SetComponent(component)
	}

	return nil
}
