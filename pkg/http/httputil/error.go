package httputil

import (
	"encoding/json"
	"net/http"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

type ErrorResponse struct {
	Message string `json:"message" example:"error message"`
}

func WriteError(rw http.ResponseWriter, msg string, code int) {
	data, err := json.Marshal(ErrorResponse{Message: msg})
	if err != nil {
		http.Error(rw, msg, code)
		return
	}

	http.Error(rw, string(data), code)
}

func WriteHTTPErrorResponse(rw http.ResponseWriter, err error) {
	switch {
	case errors.IsAlreadyExistsError(err), errors.IsInvalidStateError(err):
		WriteError(rw, err.Error(), http.StatusConflict)
	case errors.IsNotFoundError(err):
		WriteError(rw, err.Error(), http.StatusNotFound)
	case errors.IsInvalidAuthenticationError(err), errors.IsUnauthorizedError(err):
		WriteError(rw, err.Error(), http.StatusUnauthorized)
	case errors.IsInvalidParameterError(err):
		WriteError(rw, err.Error(), http.StatusUnprocessableEntity)
	case err != nil:
		WriteError(rw, "Internal server error. Please ask an admin for help or try again later", http.StatusInternalServerError)
	}
}
