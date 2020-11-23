package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/configwatcher/provider"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/config/runtime"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/httputil"
)

// InfosResponse  information exposed by the API handler.
type InfosResponse struct {
	Routers     map[string]*RouterInfoResponse     `json:"routers,omitempty"`
	Middlewares map[string]*MiddlewareInfoResponse `json:"middlewares,omitempty"`
	Services    map[string]*ServiceInfoResponse    `json:"services,omitempty"`
}

type RouterInfoResponse struct {
	*runtime.RouterInfo
	Name     string `json:"name,omitempty"`
	Provider string `json:"provider,omitempty"`
}

func NewRouterInfoResponse(name string, info *runtime.RouterInfo) *RouterInfoResponse {
	return &RouterInfoResponse{
		RouterInfo: info,
		Name:       name,
		Provider:   provider.GetName(name),
	}
}

type ServiceInfoResponse struct {
	*runtime.ServiceInfo
	ServerStatus map[string]string `json:"serverStatus,omitempty"`
	Name         string            `json:"name,omitempty"`
	Provider     string            `json:"provider,omitempty"`
	Type         string            `json:"type,omitempty"`
}

func NewServiceInfoResponse(name string, info *runtime.ServiceInfo) *ServiceInfoResponse {
	return &ServiceInfoResponse{
		ServiceInfo:  info,
		Name:         name,
		Provider:     provider.GetName(name),
		ServerStatus: info.GetAllStatus(),
		Type:         strings.ToLower(info.Service.Type()),
	}
}

type MiddlewareInfoResponse struct {
	*runtime.MiddlewareInfo
	Name     string `json:"name,omitempty"`
	Provider string `json:"provider,omitempty"`
	Type     string `json:"type,omitempty"`
}

func NewMiddlewareInfoResponse(name string, info *runtime.MiddlewareInfo) *MiddlewareInfoResponse {
	return &MiddlewareInfoResponse{
		MiddlewareInfo: info,
		Name:           name,
		Provider:       provider.GetName(name),
		Type:           strings.ToLower(info.Middleware.Type()),
	}
}

type HTTP struct {
	infos *runtime.Infos
}

func NewHTTP(infos *runtime.Infos) *HTTP {
	return &HTTP{
		infos: infos,
	}
}

func (h *HTTP) Append(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/rawdata").HandlerFunc(h.ServeHTTPInfos)
	router.Methods(http.MethodGet).Path("/http/routers").HandlerFunc(h.ServeHTTPGetRouters)
	router.Methods(http.MethodGet).Path("/http/routers/{ID}").HandlerFunc(h.ServeHTTPGetRouter)
	router.Methods(http.MethodGet).Path("/http/services").HandlerFunc(h.ServeHTTPGetServices)
	router.Methods(http.MethodGet).Path("/http/services/{ID}").HandlerFunc(h.ServeHTTPGetService)
	router.Methods(http.MethodGet).Path("/http/middlewares").HandlerFunc(h.ServeHTTPGetMiddlewares)
	router.Methods(http.MethodGet).Path("/http/middlewares/{ID}").HandlerFunc(h.ServeHTTPGetMiddleware)
}

func (h *HTTP) ServeHTTPInfos(rw http.ResponseWriter, request *http.Request) {
	result := &InfosResponse{
		Routers:     make(map[string]*RouterInfoResponse, len(h.infos.Routers)),
		Middlewares: make(map[string]*MiddlewareInfoResponse, len(h.infos.Middlewares)),
		Services:    make(map[string]*ServiceInfoResponse, len(h.infos.Services)),
	}

	for k, v := range h.infos.Routers {
		result.Routers[k] = NewRouterInfoResponse(k, v)
	}

	for k, v := range h.infos.Middlewares {
		result.Middlewares[k] = NewMiddlewareInfoResponse(k, v)
	}

	for k, v := range h.infos.Services {
		result.Services[k] = NewServiceInfoResponse(k, v)
	}

	rw.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(rw).Encode(result)
	if err != nil {
		log.FromContext(request.Context()).Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func (h *HTTP) ServeHTTPGetRouters(rw http.ResponseWriter, request *http.Request) {
	results := make([]*RouterInfoResponse, 0, len(h.infos.Routers))

	criterion := newSearchCriterion(request.URL.Query())

	for name, rt := range h.infos.Routers {
		if keepRouter(name, rt, criterion) {
			results = append(results, NewRouterInfoResponse(name, rt))
		}
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

func (h *HTTP) ServeHTTPGetRouter(rw http.ResponseWriter, request *http.Request) {
	routerID := mux.Vars(request)["ID"]

	rw.Header().Set("Content-Type", "application/json")

	router, ok := h.infos.Routers[routerID]
	if !ok {
		httputil.WriteError(rw, fmt.Sprintf("router not found: %s", routerID), http.StatusNotFound)
		return
	}

	result := NewRouterInfoResponse(routerID, router)

	err := json.NewEncoder(rw).Encode(result)
	if err != nil {
		log.FromContext(request.Context()).Error(err)
		httputil.WriteError(rw, err.Error(), http.StatusInternalServerError)
	}
}

func (h *HTTP) ServeHTTPGetServices(rw http.ResponseWriter, request *http.Request) {
	results := make([]*ServiceInfoResponse, 0, len(h.infos.Services))

	criterion := newSearchCriterion(request.URL.Query())

	for name, si := range h.infos.Services {
		if keepService(name, si, criterion) {
			results = append(results, NewServiceInfoResponse(name, si))
		}
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

func (h *HTTP) ServeHTTPGetService(rw http.ResponseWriter, request *http.Request) {
	serviceID := mux.Vars(request)["ID"]

	rw.Header().Add("Content-Type", "application/json")

	service, ok := h.infos.Services[serviceID]
	if !ok {
		httputil.WriteError(rw, fmt.Sprintf("service not found: %s", serviceID), http.StatusNotFound)
		return
	}

	result := NewServiceInfoResponse(serviceID, service)

	err := json.NewEncoder(rw).Encode(result)
	if err != nil {
		log.FromContext(request.Context()).Error(err)
		httputil.WriteError(rw, err.Error(), http.StatusInternalServerError)
	}
}

func (h *HTTP) ServeHTTPGetMiddlewares(rw http.ResponseWriter, request *http.Request) {
	results := make([]*MiddlewareInfoResponse, 0, len(h.infos.Middlewares))

	criterion := newSearchCriterion(request.URL.Query())

	for name, mi := range h.infos.Middlewares {
		if keepMiddleware(name, mi, criterion) {
			results = append(results, NewMiddlewareInfoResponse(name, mi))
		}
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

func (h *HTTP) ServeHTTPGetMiddleware(rw http.ResponseWriter, request *http.Request) {
	middlewareID := mux.Vars(request)["ID"]

	rw.Header().Set("Content-Type", "application/json")

	middleware, ok := h.infos.Middlewares[middlewareID]
	if !ok {
		httputil.WriteError(rw, fmt.Sprintf("middleware not found: %s", middlewareID), http.StatusNotFound)
		return
	}

	result := NewMiddlewareInfoResponse(middlewareID, middleware)

	err := json.NewEncoder(rw).Encode(result)
	if err != nil {
		log.FromContext(request.Context()).Error(err)
		httputil.WriteError(rw, err.Error(), http.StatusInternalServerError)
	}
}

func keepRouter(name string, item *runtime.RouterInfo, criterion *searchCriterion) bool {
	if criterion == nil {
		return true
	}

	return criterion.withStatus(item.Status) && criterion.searchIn(item.Rule, name)
}

func keepService(name string, item *runtime.ServiceInfo, criterion *searchCriterion) bool {
	if criterion == nil {
		return true
	}

	return criterion.withStatus(item.Status) && criterion.searchIn(name)
}

func keepMiddleware(name string, item *runtime.MiddlewareInfo, criterion *searchCriterion) bool {
	if criterion == nil {
		return true
	}

	return criterion.withStatus(item.Status) && criterion.searchIn(name)
}
