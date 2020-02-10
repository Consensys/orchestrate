package proto

import (
	"github.com/golang/protobuf/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

// Marshal a proto into a message assumed to be an Builder
func Marshal(pb proto.Message) ([]byte, error) {
	buf, err := proto.Marshal(pb)
	if err != nil {
		return nil, errors.EncodingError(err.Error()).SetComponent(component)
	}
	return buf, nil
}
