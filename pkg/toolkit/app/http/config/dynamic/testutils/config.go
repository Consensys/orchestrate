package testutils

import (
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/config/dynamic"
	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
)

var Configs = map[string]interface{}{
	"provider1": &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers: map[string]*dynamic.Router{
				"router-proxy": {
					Router: &traefikdynamic.Router{
						EntryPoints: []string{"ep-foo", "ep-bar"},
						Middlewares: []string{"middleware-foo"},
						Rule:        "Host(`proxy.com`)",
						Service:     "proxy",
					},
				},
			},
			Middlewares: map[string]*dynamic.Middleware{
				"middleware-foo": {
					Middleware: &traefikdynamic.Middleware{
						AddPrefix: &traefikdynamic.AddPrefix{
							Prefix: "/foo",
						},
					},
				},
				"middleware-bar": {
					Mock: &dynamic.Mock{},
				},
			},
			Services: map[string]*dynamic.Service{
				"proxy": {
					ReverseProxy: &dynamic.ReverseProxy{},
				},
				"dashboard": {
					Dashboard: &dynamic.Dashboard{},
				},
			},
		},
	},
	"provider2": &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers: map[string]*dynamic.Router{
				"router-dashboard": {
					Router: &traefikdynamic.Router{
						EntryPoints: []string{"ep-foo"},
						Middlewares: []string{"accesslog", "middleware-bar@provider1"},
						Rule:        "Host(`dashboard.com`)",
						Service:     "dashboard@provider1",
					},
				},
			},
			Middlewares: map[string]*dynamic.Middleware{
				"accesslog": {
					AccessLog: &dynamic.AccessLog{},
				},
			},
		},
	},
}

var Config = &dynamic.Configuration{
	HTTP: &dynamic.HTTPConfiguration{
		Routers: map[string]*dynamic.Router{
			"router-proxy@provider1": {
				Router: &traefikdynamic.Router{
					EntryPoints: []string{"ep-foo", "ep-bar"},
					Middlewares: []string{"middleware-foo@provider1"},
					Rule:        "Host(`proxy.com`)",
					Service:     "proxy@provider1",
				},
			},
			"router-dashboard@provider2": {
				Router: &traefikdynamic.Router{
					EntryPoints: []string{"ep-foo"},
					Middlewares: []string{"accesslog@provider2", "middleware-bar@provider1"},
					Rule:        "Host(`dashboard.com`)",
					Service:     "dashboard@provider1",
				},
			},
		},
		Middlewares: map[string]*dynamic.Middleware{
			"middleware-foo@provider1": {
				Middleware: &traefikdynamic.Middleware{
					AddPrefix: &traefikdynamic.AddPrefix{
						Prefix: "/foo",
					},
				},
			},
			"middleware-bar@provider1": {
				Mock: &dynamic.Mock{},
			},
			"accesslog@provider2": {
				AccessLog: &dynamic.AccessLog{},
			},
		},
		Services: map[string]*dynamic.Service{
			"proxy@provider1": {
				ReverseProxy: &dynamic.ReverseProxy{},
			},
			"dashboard@provider1": {
				Dashboard: &dynamic.Dashboard{},
			},
		},
	},
}
