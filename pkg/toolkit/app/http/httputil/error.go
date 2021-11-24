package httputil

import (
	"encoding/json"
	"net/http"

	"github.com/consensys/orchestrate/pkg/errors"
)

var internalErrMsg = "Internal server error. Please ask an admin for help or try again later"
var internalDepErrMsg = "Failed dependency. Please ask an admin for help or try again later"

type ErrorResponse struct {
	Message string `json:"message" example:"error message"`
	Code    uint64 `json:"code,omitempty" example:"24000"`
}

// @deprecated: Migrate every usage of it to WriteHTTPErrorResponse
func WriteError(rw http.ResponseWriter, msg string, code int) {
	var err error
	switch code {
	case http.StatusBadRequest:
		err = errors.InvalidFormatError(msg)
	case http.StatusConflict:
		err = errors.ConflictedError(msg)
	case http.StatusUnprocessableEntity:
		err = errors.InvalidParameterError(msg)
	case http.StatusNotFound:
		err = errors.NotFoundError(msg)
	default:
		err = errors.InternalError(msg)
	}

	writeErrorResponse(rw, code, err)
}

func WriteHTTPErrorResponse(rw http.ResponseWriter, err error) {
	switch {
	case errors.IsAlreadyExistsError(err), errors.IsInvalidStateError(err):
		writeErrorResponse(rw, http.StatusConflict, err)
	case errors.IsNotFoundError(err):
		writeErrorResponse(rw, http.StatusNotFound, err)
	case errors.IsInvalidAuthenticationError(err), errors.IsUnauthorizedError(err):
		writeErrorResponse(rw, http.StatusUnauthorized, err)
	case errors.IsInvalidFormatError(err):
		writeErrorResponse(rw, http.StatusBadRequest, err)
	case errors.IsInvalidParameterError(err), errors.IsEncodingError(err):
		writeErrorResponse(rw, http.StatusUnprocessableEntity, err)
	case errors.IsPostgresConnectionError(err), errors.IsKafkaConnectionError(err):
		writeErrorResponse(rw, http.StatusFailedDependency, errors.FromError(err).SetMessage(internalDepErrMsg))
	case err != nil:
		writeErrorResponse(rw, http.StatusInternalServerError, errors.FromError(err).SetMessage(internalErrMsg))
	}
}

func writeErrorResponse(rw http.ResponseWriter, status int, err error) {
	msg, e := json.Marshal(ErrorResponse{
		Message: errors.FromError(err).SetComponent("").Error(),
		Code:    errors.FromError(err).GetCode(),
	})
	if e != nil {
		http.Error(rw, e.Error(), status)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("X-Content-Type-Options", "nosniff")
	rw.WriteHeader(status)
	_, _ = rw.Write(msg)
}
