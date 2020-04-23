package json

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
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
		errMessage := "failed to decode request body"
		log.WithError(err).Error(errMessage)
		return errors.InvalidFormatError(err.Error()).ExtendComponent(component)
	}

	err = utils.GetValidator().Struct(req)
	if err != nil {
		if ves, ok := err.(validator.ValidationErrors); ok {
			var errMessage string
			for _, fe := range ves {
				errMessage += fmt.Sprintf(" field validation for '%s' failed on the '%s' tag", fe.Field(), fe.Tag())
			}

			log.WithError(err).Error(errMessage)
			return errors.InvalidParameterError("invalid body, with:%s", errMessage).ExtendComponent(component)
		}

		errMessage := "Invalid body"
		log.WithError(err).Error(errMessage)
		return errors.InvalidFormatError(errMessage).ExtendComponent(component)
	}

	return nil
}
