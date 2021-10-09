package dynamic

import (
	traefikdynamic "github.com/traefik/traefik/v2/pkg/config/dynamic"
)

// +k8s:deepcopy-gen=true

// HTTPConfiguration contains all the HTTP configuration parameters.
type HTTPConfiguration struct {
	Routers     map[string]*Router     `json:"routers,omitempty" toml:"routers,omitempty" yaml:"routers,omitempty"`
	Middlewares map[string]*Middleware `json:"middlewares,omitempty" toml:"middlewares,omitempty" yaml:"middlewares,omitempty"`
	Services    map[string]*Service    `json:"services,omitempty" toml:"services,omitempty" yaml:"services,omitempty"`
}

func FromTraefikHTTPConfig(traefikConf *traefikdynamic.HTTPConfiguration) *HTTPConfiguration {
	if traefikConf == nil {
		return nil
	}

	conf := &HTTPConfiguration{}

	if len(traefikConf.Routers) > 0 {
		conf.Routers = make(map[string]*Router)
		for k, router := range traefikConf.Routers {
			conf.Routers[k] = FromTraefikRouter(router)
		}
	}

	if len(traefikConf.Middlewares) > 0 {
		conf.Middlewares = make(map[string]*Middleware)
		for k, middleware := range traefikConf.Middlewares {
			conf.Middlewares[k] = FromTraefikMiddleware(middleware)
		}
	}

	if len(traefikConf.Services) > 0 {
		conf.Services = make(map[string]*Service)
		for k, service := range traefikConf.Services {
			conf.Services[k] = FromTraefikService(service)
		}
	}

	return conf
}

func ToTraefikHTTPConfig(conf *HTTPConfiguration) *traefikdynamic.HTTPConfiguration {
	if conf == nil {
		return nil
	}

	traefikConf := &traefikdynamic.HTTPConfiguration{}

	if len(conf.Routers) > 0 {
		traefikConf.Routers = make(map[string]*traefikdynamic.Router)
		for k, router := range conf.Routers {
			traefikConf.Routers[k] = ToTraefikRouter(router)
		}
	}

	if len(conf.Middlewares) > 0 {
		traefikConf.Middlewares = make(map[string]*traefikdynamic.Middleware)
		for k, middleware := range conf.Middlewares {
			traefikConf.Middlewares[k] = ToTraefikMiddleware(middleware)
		}
	}

	if len(conf.Services) > 0 {
		traefikConf.Services = make(map[string]*traefikdynamic.Service)
		for k, service := range conf.Services {
			traefikConf.Services[k] = ToTraefikService(service)
		}
	}

	return traefikConf
}
