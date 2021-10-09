package dynamic

import (
	traefikdynamic "github.com/traefik/traefik/v2/pkg/config/dynamic"
)

// +k8s:deepcopy-gen=true

type Router struct {
	*traefikdynamic.Router
}

// func (rt *Router) DeepCopy() *Router {
// 	return &Router{
// 		Router: &traefikdynamic.Router{
// 			EntryPoints: rt.EntryPoints,
// 			Middlewares: rt.Middlewares,
// 			Service:     rt.Service,
// 			Rule:        rt.Rule,
// 			Priority:    rt.Priority,
// 			TLS:         rt.TLS,
// 		},
// 	}
// }

func FromTraefikRouter(router *traefikdynamic.Router) *Router {
	return &Router{
		Router: router,
	}
}

func ToTraefikRouter(router *Router) *traefikdynamic.Router {
	return router.Router
}
