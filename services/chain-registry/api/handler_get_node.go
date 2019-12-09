package api

import (
	"encoding/json"
	"net/http"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/gorilla/mux"
)

type nodeRepresentation struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

func (h Handler) getNode(rw http.ResponseWriter, request *http.Request) {
	log.FromContext(request.Context()).Infof("getNode")

	nodeID := mux.Vars(request)["nodeID"]

	rw.Header().Set("Content-Type", "application/json")

	result := nodeRepresentation{
		ID:   nodeID,
		Name: "test-node",
		URL:  "localhost:8545",
	}

	err := json.NewEncoder(rw).Encode(result)
	if err != nil {
		log.FromContext(request.Context()).Error(err)
		writeError(rw, err.Error(), http.StatusInternalServerError)
	}
}
