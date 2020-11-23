package httptrace

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptrace"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/config/dynamic"
)

type Builder struct{}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Build(ctx context.Context, _ string, configuration interface{}) (mid func(http.Handler) http.Handler, respModifier func(resp *http.Response) error, err error) {
	cfg, ok := configuration.(*dynamic.HTTPTrace)
	if !ok {
		return nil, nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	return New().Handler, nil, nil
}

type HTTPTrace struct{}

func New() *HTTPTrace {
	return &HTTPTrace{}
}

func (t *HTTPTrace) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		t.ServeHTTP(rw, req, h)
	})
}

func (t *HTTPTrace) ServeHTTP(rw http.ResponseWriter, req *http.Request, next http.Handler) {
	tr := NewLoggerTrace(log.FromContext(req.Context()))
	ctx := httptrace.WithClientTrace(req.Context(), tr)
	next.ServeHTTP(rw, req.WithContext(ctx))
}

func NewLoggerTrace(logger logrus.FieldLogger) *httptrace.ClientTrace {
	return &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			logger.
				WithField("host", hostPort).
				Info("GetConn")
		},
		GotConn: func(info httptrace.GotConnInfo) {
			logger.
				WithField("reused", info.Reused).
				WithField("was-idle", info.WasIdle).
				WithField("idletime", info.IdleTime).
				Info("GotConn")
		},
		WroteHeaderField: func(key string, value []string) {
			logger.
				WithField(key, value).
				Info("WroteHeaderField")
		},
	}
}
