package utils

import (
	"encoding/json"
	"html"
	"net/http"
	"net/url"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

type APIError struct {
	Message string `json:"message"`
}

func ToFilters(values url.Values) map[string]string {
	filters := make(map[string]string)
	for key := range values {
		k := html.EscapeString(key)
		v := html.EscapeString(values.Get(key))
		if k != "" && v != "" {
			filters[k] = v
		}
	}
	return filters
}

func HandleStoreError(rw http.ResponseWriter, err error) {
	switch {
	case errors.IsAlreadyExistsError(err):
		WriteError(rw, err.Error(), http.StatusConflict)
	case errors.IsNotFoundError(err):
		WriteError(rw, err.Error(), http.StatusNotFound)
	case errors.IsDataError(err):
		WriteError(rw, err.Error(), http.StatusBadRequest)
	case err != nil:
		WriteError(rw, err.Error(), http.StatusInternalServerError)
	}
}

func WriteError(rw http.ResponseWriter, msg string, code int) {
	data, _ := json.Marshal(APIError{Message: msg})
	http.Error(rw, string(data), code)
}
