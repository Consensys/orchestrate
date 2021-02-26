package router

import (
	"context"
	"crypto/tls"
	"net/http"
)

//go:generate mockgen -source=builder.go -destination=mock/mock.go -package=mock

type Builder interface {
	Build(ctx context.Context, entryPointNames []string, configuration interface{}) (map[string]*Router, error)
}

type Router struct {
	HTTP           http.Handler
	HTTPS          http.Handler
	TLSConfig      *tls.Config
	HostTLSConfigs map[string]*tls.Config
}
