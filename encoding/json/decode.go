package json

import (
	"encoding/json"

	errors "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
)

// Unmarshal parses the JSON-encoded data and stores the result in the value pointed to by v
func Unmarshal(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return errors.EncodingError(err).SetComponent(component)
	}
	return nil
}
