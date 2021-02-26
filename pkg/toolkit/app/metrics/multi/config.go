package multi

import (
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/http/config/dynamic"
)

// dynamicConfig holds the current set of routers, and services
// corresponding to the current state of an HTTP router
// It provides a performant way to check whether the collected metrics belong to the
// current configuration or to an outdated one.
type DynamicConfig struct {
	routers  map[string]bool
	services map[string]map[string]bool
}

func NewDynamicConfig(conf *dynamic.Configuration) *DynamicConfig {
	cfg := &DynamicConfig{
		routers:  make(map[string]bool),
		services: make(map[string]map[string]bool),
	}

	if conf == nil || conf.HTTP == nil {
		return cfg
	}

	for rtName := range conf.HTTP.Routers {
		cfg.routers[rtName] = true
	}

	for serviceName, service := range conf.HTTP.Services {
		cfg.services[serviceName] = make(map[string]bool)
		if service.ReverseProxy != nil {
			for _, server := range service.ReverseProxy.LoadBalancer.Servers {
				cfg.services[serviceName][server.URL] = true
			}
		}
	}

	return cfg
}

func (cfg *DynamicConfig) hasService(serviceName string) bool {
	_, ok := cfg.services[serviceName]
	return ok
}

func (cfg *DynamicConfig) hasServerURL(serviceName, serverURL string) bool {
	if service, hasService := cfg.services[serviceName]; hasService {
		_, ok := service[serverURL]
		return ok
	}
	return false
}
