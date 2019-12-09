package api

import (
	"encoding/json"
	"net/http"

	"github.com/containous/traefik/v2/pkg/config/runtime"
	"github.com/gorilla/mux"
)

type apiError struct {
	Message string `json:"message"`
}

func writeError(rw http.ResponseWriter, msg string, code int) {
	data, err := json.Marshal(apiError{Message: msg})
	if err != nil {
		http.Error(rw, msg, code)
		return
	}

	http.Error(rw, string(data), code)
}

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

// Add internal routes to router
func (h *Handler) Append(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/api/nodes/{nodeID}").HandlerFunc(h.getNode)
}

type Builder func(config *runtime.Configuration) http.Handler
