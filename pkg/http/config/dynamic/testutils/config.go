package testutils

import (
	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
)

var Configs = map[string]interface{}{
	"provider1": &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers: map[string]*dynamic.Router{
				"router-proxy": &dynamic.Router{
					Router: &traefikdynamic.Router{
						EntryPoints: []string{"ep-foo", "ep-bar"},
						Middlewares: []string{"middleware-foo"},
						Rule:        "Host(`proxy.com`)",
						Service:     "proxy",
					},
				},
			},
			Middlewares: map[string]*dynamic.Middleware{
				"middleware-foo": &dynamic.Middleware{
					Middleware: &traefikdynamic.Middleware{
						AddPrefix: &traefikdynamic.AddPrefix{
							Prefix: "/foo",
						},
					},
				},
				"middleware-bar": &dynamic.Middleware{
					Mock: &dynamic.Mock{},
				},
			},
			Services: map[string]*dynamic.Service{
				"proxy": &dynamic.Service{
					ReverseProxy: &dynamic.ReverseProxy{},
				},
				"dashboard": &dynamic.Service{
					Dashboard: &dynamic.Dashboard{},
				},
			},
		},
	},
	"provider2": &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers: map[string]*dynamic.Router{
				"router-dashboard": &dynamic.Router{
					Router: &traefikdynamic.Router{
						EntryPoints: []string{"ep-foo"},
						Middlewares: []string{"accesslog", "middleware-bar@provider1"},
						Rule:        "Host(`dashboard.com`)",
						Service:     "dashboard@provider1",
					},
				},
			},
			Middlewares: map[string]*dynamic.Middleware{
				"accesslog": &dynamic.Middleware{
					AccessLog: &dynamic.AccessLog{},
				},
			},
		},
	},
}

var Config = &dynamic.Configuration{
	HTTP: &dynamic.HTTPConfiguration{
		Routers: map[string]*dynamic.Router{
			"router-proxy@provider1": &dynamic.Router{
				Router: &traefikdynamic.Router{
					EntryPoints: []string{"ep-foo", "ep-bar"},
					Middlewares: []string{"middleware-foo@provider1"},
					Rule:        "Host(`proxy.com`)",
					Service:     "proxy@provider1",
				},
			},
			"router-dashboard@provider2": &dynamic.Router{
				Router: &traefikdynamic.Router{
					EntryPoints: []string{"ep-foo"},
					Middlewares: []string{"accesslog@provider2", "middleware-bar@provider1"},
					Rule:        "Host(`dashboard.com`)",
					Service:     "dashboard@provider1",
				},
			},
		},
		Middlewares: map[string]*dynamic.Middleware{
			"middleware-foo@provider1": &dynamic.Middleware{
				Middleware: &traefikdynamic.Middleware{
					AddPrefix: &traefikdynamic.AddPrefix{
						Prefix: "/foo",
					},
				},
			},
			"middleware-bar@provider1": &dynamic.Middleware{
				Mock: &dynamic.Mock{},
			},
			"accesslog@provider2": &dynamic.Middleware{
				AccessLog: &dynamic.AccessLog{},
			},
		},
		Services: map[string]*dynamic.Service{
			"proxy@provider1": &dynamic.Service{
				ReverseProxy: &dynamic.ReverseProxy{},
			},
			"dashboard@provider1": &dynamic.Service{
				Dashboard: &dynamic.Dashboard{},
			},
		},
	},
}
