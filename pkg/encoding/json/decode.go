package json

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/ConsenSys/orchestrate/pkg/utils"
	"github.com/go-playground/validator/v10"

	"github.com/ConsenSys/orchestrate/pkg/errors"
)

// Unmarshal parses the JSON-encoded data and stores the result in the value pointed to by v
func Unmarshal(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return errors.EncodingError(err.Error()).SetComponent(component)
	}
	return nil
}

func UnmarshalBody(body io.Reader, req interface{}) error {
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields() // Force errors if unknown fields
	err := dec.Decode(req)
	if err != nil {
		return errors.InvalidFormatError("failed to decode request body").AppendReason(err.Error())
	}

	err = utils.GetValidator().Struct(req)
	if err != nil {
		if ves, ok := err.(validator.ValidationErrors); ok {
			var errMessage string
			for _, fe := range ves {
				errMessage += fmt.Sprintf("field validation for '%s' failed on the '%s' tag", fe.Field(), fe.Tag())
			}

			return errors.InvalidParameterError("invalid body").AppendReason(errMessage)
		}

		return errors.InvalidFormatError("invalid body")
	}

	return nil
}
