package json

import (
	"encoding/json"

	"github.com/consensys/orchestrate/pkg/errors"
)

// Marshal returns the JSON encoding of v
func Marshal(v interface{}) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, errors.EncodingError(err.Error()).SetComponent(component)
	}
	return b, nil
}
