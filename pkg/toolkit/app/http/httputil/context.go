package httputil

import (
	"context"
)

type httpCtxKey string

const (
	epNameKey         httpCtxKey = "entrypoint"
	routerNameKey     httpCtxKey = "router"
	serviceNameKey    httpCtxKey = "service"
	middlewareNameKey httpCtxKey = "middleware"
)

func WithEntryPoint(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, epNameKey, name)
}

func EntryPointFromContext(ctx context.Context) string {
	v, ok := ctx.Value(epNameKey).(string)
	if ok {
		return v
	}
	return ""
}

func WithRouter(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, routerNameKey, name)
}

func RouterFromContext(ctx context.Context) string {
	v, ok := ctx.Value(routerNameKey).(string)
	if ok {
		return v
	}
	return ""
}

func WithService(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, serviceNameKey, name)
}

func ServiceFromContext(ctx context.Context) string {
	v, ok := ctx.Value(serviceNameKey).(string)
	if ok {
		return v
	}
	return ""
}

func WithMiddleware(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, middlewareNameKey, name)
}

func MiddlewareFromContext(ctx context.Context) string {
	v, ok := ctx.Value(middlewareNameKey).(string)
	if ok {
		return v
	}
	return ""
}
