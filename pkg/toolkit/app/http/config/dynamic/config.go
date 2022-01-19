package dynamic

import (
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/configwatcher/provider"
	traefikdynamic "github.com/traefik/traefik/v2/pkg/config/dynamic"
	"github.com/traefik/traefik/v2/pkg/log"
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
