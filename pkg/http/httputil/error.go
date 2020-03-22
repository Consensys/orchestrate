package httputil

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func WriteError(rw http.ResponseWriter, msg string, code int) {
	data, err := json.Marshal(ErrorResponse{Message: msg})
	if err != nil {
		http.Error(rw, msg, code)
		return
	}

	http.Error(rw, string(data), code)
}
