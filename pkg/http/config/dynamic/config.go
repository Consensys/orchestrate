package dynamic

import (
	"reflect"

	"github.com/ConsenSys/orchestrate/pkg/configwatcher/provider"
	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/log"
)

// Configuration is the root of the dynamic configuration
type Configuration struct {
	HTTP *HTTPConfiguration               `json:"http,omitempty" toml:"http,omitempty" yaml:"http,omitempty"`
	TLS  *traefikdynamic.TLSConfiguration `json:"tls,omitempty" toml:"tls,omitempty" yaml:"tls,omitempty"`
}

func NewConfig() *Configuration {
	return &Configuration{
		HTTP: &HTTPConfiguration{
			Routers:     make(map[string]*Router),
			Middlewares: make(map[string]*Middleware),
			Services:    make(map[string]*Service),
		},
	}
}

func FromTraefikConfig(traefikConf *traefikdynamic.Configuration) *Configuration {
	if traefikConf == nil {
		return nil
	}

	return &Configuration{
		HTTP: FromTraefikHTTPConfig(traefikConf.HTTP),
		TLS:  traefikConf.TLS,
	}
}

func ToTraefikConfig(conf *Configuration) *traefikdynamic.Configuration {
	if conf == nil {
		return nil
	}

	return &traefikdynamic.Configuration{
		HTTP: ToTraefikHTTPConfig(conf.HTTP),
		TLS:  conf.TLS,
	}
}

// Merge Merges multiple configurations.
func Merge(configurations map[string]interface{}) interface{} {
	logger := log.WithoutContext()

	configuration := NewConfig()
	for providerName, cfg := range configurations {
		conf, ok := cfg.(*Configuration)
		if !ok {
			logger.Errorf("Found invalid configuration type while merging (expected %T but got %T)", conf, cfg)
			return nil
		}

		if conf.HTTP != nil {
			for serviceName, service := range conf.HTTP.Services {
				srv := service.DeepCopy()
				configuration.HTTP.Services[provider.QualifyName(providerName, serviceName)] = srv
			}

			for middlewareName, middleware := range conf.HTTP.Middlewares {
				mid := middleware.DeepCopy()
				configuration.HTTP.Middlewares[provider.QualifyName(providerName, middlewareName)] = mid
			}

			for routerName, router := range conf.HTTP.Routers {
				rt := router.DeepCopy()
				var midNames []string
				for _, midName := range router.Middlewares {
					midNames = append(midNames, provider.QualifyName(providerName, midName))
				}
				rt.Middlewares = midNames

				rt.Service = provider.QualifyName(providerName, router.Service)

				configuration.HTTP.Routers[provider.QualifyName(providerName, routerName)] = rt
			}
		}
	}

	return configuration
}

// AddService Adds a service to a configurations.
func AddService(configuration *HTTPConfiguration, serviceName string, service *Service) bool {
	if _, ok := configuration.Services[serviceName]; !ok {
		configuration.Services[serviceName] = service
		return true
	}

	if configuration.Services[serviceName].ReverseProxy == nil || service.ReverseProxy == nil || !configuration.Services[serviceName].ReverseProxy.Mergeable(service.ReverseProxy) {
		return false
	}

	configuration.Services[serviceName].ReverseProxy.LoadBalancer.Servers = append(configuration.Services[serviceName].ReverseProxy.LoadBalancer.Servers, service.ReverseProxy.LoadBalancer.Servers...)
	return true
}

// AddRouter Adds a router to a configurations.
func AddRouter(configuration *HTTPConfiguration, routerName string, router *Router) bool {
	if _, ok := configuration.Routers[routerName]; !ok {
		configuration.Routers[routerName] = router
		return true
	}

	return reflect.DeepEqual(configuration.Routers[routerName], router)
}

// AddMiddleware Adds a middleware to a configurations.
func AddMiddleware(configuration *HTTPConfiguration, middlewareName string, middleware *Middleware) bool {
	if _, ok := configuration.Middlewares[middlewareName]; !ok {
		configuration.Middlewares[middlewareName] = middleware
		return true
	}

	return reflect.DeepEqual(configuration.Middlewares[middlewareName], middleware)
}
