package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/consensys/orchestrate/pkg/toolkit/app/http/httputil"
	"github.com/gorilla/mux"
	traefikstatic "github.com/traefik/traefik/v2/pkg/config/static"
	"github.com/traefik/traefik/v2/pkg/log"
)

type EntryPointResponse struct {
	*traefikstatic.EntryPoint
	Name string `json:"name,omitempty"`
}

type EntryPoint struct {
	staticCfg *traefikstatic.Configuration
}

func NewEntryPoint(staticCfg *traefikstatic.Configuration) *EntryPoint {
	return &EntryPoint{
		staticCfg: staticCfg,
	}
}

func (h *EntryPoint) Append(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/entrypoints").HandlerFunc(h.ServeHTTPGetEntryPoints)
	router.Methods(http.MethodGet).Path("/entrypoints/{ID}").HandlerFunc(h.ServeHTTPGetEntryPoint)
}

func (h *EntryPoint) ServeHTTPGetEntryPoints(rw http.ResponseWriter, request *http.Request) {
	results := make([]*EntryPointResponse, 0, len(h.staticCfg.EntryPoints))

	for name, ep := range h.staticCfg.EntryPoints {
		results = append(results, &EntryPointResponse{
			EntryPoint: ep,
			Name:       name,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	rw.Header().Set("Content-Type", "application/json")

	pageInfo, err := pagination(request, len(results))
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.Header().Set(nextPageHeader, strconv.Itoa(pageInfo.nextPage))

	err = json.NewEncoder(rw).Encode(results[pageInfo.startIndex:pageInfo.endIndex])
	if err != nil {
		log.FromContext(request.Context()).Error(err)
		httputil.WriteError(rw, err.Error(), http.StatusInternalServerError)
	}
}

func (h *EntryPoint) ServeHTTPGetEntryPoint(rw http.ResponseWriter, request *http.Request) {
	entryPointID := mux.Vars(request)["ID"]

	rw.Header().Set("Content-Type", "application/json")

	ep, ok := h.staticCfg.EntryPoints[entryPointID]
	if !ok {
		httputil.WriteError(rw, fmt.Sprintf("entry point not found: %s", entryPointID), http.StatusNotFound)
		return
	}

	result := &EntryPointResponse{
		EntryPoint: ep,
		Name:       entryPointID,
	}

	err := json.NewEncoder(rw).Encode(result)
	if err != nil {
		log.FromContext(request.Context()).Error(err)
		httputil.WriteError(rw, err.Error(), http.StatusInternalServerError)
	}
}
