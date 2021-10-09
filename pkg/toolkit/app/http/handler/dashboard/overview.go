package dashboard

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/consensys/orchestrate/pkg/toolkit/app/http/config/runtime"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/httputil"
	"github.com/gorilla/mux"
	traefikstatic "github.com/traefik/traefik/v2/pkg/config/static"
	"github.com/traefik/traefik/v2/pkg/log"
)

type OverviewResponse struct {
	HTTP      *HTTPOverview     `json:"http"`
	Features  *FeaturesResponse `json:"features,omitempty"`
	Providers []string          `json:"providers,omitempty"`
}

type HTTPOverview struct {
	Routers     *SectionResponse `json:"routers,omitempty"`
	Services    *SectionResponse `json:"services,omitempty"`
	Middlewares *SectionResponse `json:"middlewares,omitempty"`
}

type SectionResponse struct {
	Total    int `json:"total"`
	Warnings int `json:"warnings"`
	Errors   int `json:"errors"`
}

type FeaturesResponse struct {
	Tracing   string `json:"tracing"`
	Metrics   string `json:"metrics"`
	AccessLog bool   `json:"accessLog"`
	// TODO add certificates resolvers
}

type Overview struct {
	staticCfg *traefikstatic.Configuration
	infos     *runtime.Infos
}

func NewOverview(staticCfg *traefikstatic.Configuration, infos *runtime.Infos) *Overview {
	return &Overview{
		staticCfg: staticCfg,
		infos:     infos,
	}
}

func (h *Overview) Append(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/overview").HandlerFunc(h.ServeHTTP)
}

func (h *Overview) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	result := &OverviewResponse{
		HTTP: &HTTPOverview{
			Routers:     getHTTPRouterSection(h.infos.Routers),
			Services:    getHTTPServiceSection(h.infos.Services),
			Middlewares: getHTTPMiddlewareSection(h.infos.Middlewares),
		},
		Features:  getFeatures(h.staticCfg),
		Providers: getProviders(h.staticCfg),
	}

	rw.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(rw).Encode(result)
	if err != nil {
		log.FromContext(request.Context()).Error(err)
		httputil.WriteError(rw, err.Error(), http.StatusInternalServerError)
	}
}

func getHTTPRouterSection(routers map[string]*runtime.RouterInfo) *SectionResponse {
	var countErrors int
	var countWarnings int
	for _, rt := range routers {
		switch rt.Status {
		case runtime.StatusDisabled:
			countErrors++
		case runtime.StatusWarning:
			countWarnings++
		}
	}

	return &SectionResponse{
		Total:    len(routers),
		Warnings: countWarnings,
		Errors:   countErrors,
	}
}

func getHTTPServiceSection(services map[string]*runtime.ServiceInfo) *SectionResponse {
	var countErrors int
	var countWarnings int
	for _, svc := range services {
		switch svc.Status {
		case runtime.StatusDisabled:
			countErrors++
		case runtime.StatusWarning:
			countWarnings++
		}
	}

	return &SectionResponse{
		Total:    len(services),
		Warnings: countWarnings,
		Errors:   countErrors,
	}
}

func getHTTPMiddlewareSection(middlewares map[string]*runtime.MiddlewareInfo) *SectionResponse {
	var countErrors int
	var countWarnings int
	for _, mid := range middlewares {
		switch mid.Status {
		case runtime.StatusDisabled:
			countErrors++
		case runtime.StatusWarning:
			countWarnings++
		}
	}

	return &SectionResponse{
		Total:    len(middlewares),
		Warnings: countWarnings,
		Errors:   countErrors,
	}
}

func getProviders(conf *traefikstatic.Configuration) []string {
	if conf.Providers == nil {
		return nil
	}

	var providers []string

	v := reflect.ValueOf(conf.Providers).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Ptr && field.Elem().Kind() == reflect.Struct {
			if !field.IsNil() {
				providers = append(providers, v.Type().Field(i).Name)
			}
		}
	}

	return providers
}

func getFeatures(conf *traefikstatic.Configuration) *FeaturesResponse {
	return &FeaturesResponse{
		Tracing:   getTracing(conf),
		Metrics:   getMetrics(conf),
		AccessLog: conf.AccessLog != nil,
	}
}

func getMetrics(conf *traefikstatic.Configuration) string {
	if conf.Metrics == nil {
		return ""
	}

	v := reflect.ValueOf(conf.Metrics).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Ptr && field.Elem().Kind() == reflect.Struct {
			if !field.IsNil() {
				return v.Type().Field(i).Name
			}
		}
	}

	return ""
}

func getTracing(conf *traefikstatic.Configuration) string {
	if conf.Tracing == nil {
		return ""
	}

	v := reflect.ValueOf(conf.Tracing).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Ptr && field.Elem().Kind() == reflect.Struct {
			if !field.IsNil() {
				return v.Type().Field(i).Name
			}
		}
	}

	return ""
}
