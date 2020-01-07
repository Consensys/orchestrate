// All scripts in server are highly inspired from Traefik server
// c.f. https://github.com/containous/traefik/tree/v2.0.5/pkg/server

package server

import (
	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/log"
	"github.com/containous/traefik/v2/pkg/tls"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/server/utils"
)

func mergeConfiguration(configurations dynamic.Configurations) dynamic.Configuration {
	conf := dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers:     make(map[string]*dynamic.Router),
			Middlewares: make(map[string]*dynamic.Middleware),
			Services:    make(map[string]*dynamic.Service),
		},
		TCP: &dynamic.TCPConfiguration{
			Routers:  make(map[string]*dynamic.TCPRouter),
			Services: make(map[string]*dynamic.TCPService),
		},
		TLS: &dynamic.TLSConfiguration{
			Stores:  make(map[string]tls.Store),
			Options: make(map[string]tls.Options),
		},
	}

	var defaultTLSOptionProviders []string
	for provider, configuration := range configurations {
		if configuration.HTTP != nil {
			for routerName, router := range configuration.HTTP.Routers {
				conf.HTTP.Routers[utils.MakeQualifiedName(provider, routerName)] = router
			}
			for middlewareName, middleware := range configuration.HTTP.Middlewares {
				conf.HTTP.Middlewares[utils.MakeQualifiedName(provider, middlewareName)] = middleware
			}
			for serviceName, service := range configuration.HTTP.Services {
				conf.HTTP.Services[utils.MakeQualifiedName(provider, serviceName)] = service
			}
		}

		if configuration.TCP != nil {
			for routerName, router := range configuration.TCP.Routers {
				conf.TCP.Routers[utils.MakeQualifiedName(provider, routerName)] = router
			}
			for serviceName, service := range configuration.TCP.Services {
				conf.TCP.Services[utils.MakeQualifiedName(provider, serviceName)] = service
			}
		}

		if configuration.TLS != nil {
			conf.TLS.Certificates = append(conf.TLS.Certificates, configuration.TLS.Certificates...)

			for key, store := range configuration.TLS.Stores {
				conf.TLS.Stores[key] = store
			}

			for tlsOptionsName, options := range configuration.TLS.Options {
				if tlsOptionsName != "default" {
					tlsOptionsName = utils.MakeQualifiedName(provider, tlsOptionsName)
				} else {
					defaultTLSOptionProviders = append(defaultTLSOptionProviders, provider)
				}

				conf.TLS.Options[tlsOptionsName] = options
			}
		}
	}

	if len(defaultTLSOptionProviders) == 0 {
		conf.TLS.Options["default"] = tls.Options{}
	} else if len(defaultTLSOptionProviders) > 1 {
		log.WithoutContext().Errorf("Default TLS Options defined multiple times in %v", defaultTLSOptionProviders)
		// We do not set an empty tls.TLS{} as above so that we actually get a "cascading failure" later on,
		// i.e. routers depending on this missing TLS option will fail to initialize as well.
		delete(conf.TLS.Options, "default")
	}

	return conf
}
