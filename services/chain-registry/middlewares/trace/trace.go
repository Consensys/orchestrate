package trace

import (
	"net/http"
	"net/http/httptrace"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/sirupsen/logrus"
)

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

type RequestTracer struct {
	next http.Handler
}

func New(next http.Handler) *RequestTracer {
	return &RequestTracer{
		next: next,
	}
}

func (t *RequestTracer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	tr := NewLoggerTrace(log.FromContext(req.Context()))
	ctx := httptrace.WithClientTrace(req.Context(), tr)
	t.next.ServeHTTP(rw, req.WithContext(ctx))
}
