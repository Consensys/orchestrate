package proto

import (
	"github.com/golang/protobuf/proto"
	errors "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
)

// Unmarshal parses the protocol buffer representation in `buf` and places the decoded result in `pb`
func Unmarshal(buf []byte, pb proto.Message) error {
	// Cast message into protobuffer
	e := proto.Unmarshal(buf, pb)
	if e != nil {
		return errors.EncodingError(e).SetComponent(component)
	}

	return nil
}
