package proxy

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	gohttputil "net/http/httputil"
	"net/url"
	"time"

	"github.com/oxtoacart/bpool"

	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"github.com/containous/traefik/v2/pkg/log"
	"github.com/containous/traefik/v2/pkg/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/httputil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/loadbalancer"
)

// StatusClientClosedRequest non-standard HTTP status code for client disconnection
const StatusClientClosedRequest = 499

// StatusClientClosedRequestText non-standard HTTP status for client disconnection
const StatusClientClosedRequestText = "Client Closed Envelope"

type Builder struct {
	transport http.RoundTripper
	bpool     gohttputil.BufferPool

	lb *loadbalancer.Builder
}

func NewBuilder(staticCfg *traefikstatic.Configuration, pool gohttputil.BufferPool) (*Builder, error) {
	t, err := NewTransport(staticCfg.ServersTransport)
	if err != nil {
		return nil, err
	}

	if pool == nil {
		pool = bpool.NewBytePool(32, 1024)
	}

	return &Builder{
		transport: t,
		bpool:     pool,
		lb:        loadbalancer.NewBuilder(),
	}, nil
}

func (b *Builder) Build(ctx context.Context, name string, configuration interface{}, respModifier func(*http.Response) error) (http.Handler, error) {
	cfg, ok := configuration.(*dynamic.ReverseProxy)
	if !ok {
		return nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	if cfg.LoadBalancer == nil {
		return New(cfg, b.transport, b.bpool, respModifier)
	}

	lb, lbRespModifier, err := b.lb.Build(ctx, name, cfg.LoadBalancer)
	if err != nil {
		return nil, err
	}

	respModifier = httputil.CombineResponseModifiers(respModifier, lbRespModifier)
	proxy, err := New(cfg, b.transport, b.bpool, respModifier)
	if err != nil {
		return nil, err
	}

	return lb(proxy), nil
}

func New(cfg *dynamic.ReverseProxy, transport http.RoundTripper, pool gohttputil.BufferPool, respModifier func(*http.Response) error) (*gohttputil.ReverseProxy, error) {
	var flushInterval types.Duration
	if cfg.ResponseForwarding != nil && cfg.ResponseForwarding.FlushInterval != "" {
		err := flushInterval.Set(cfg.ResponseForwarding.FlushInterval)
		if err != nil {
			return nil, fmt.Errorf("error creating flush interval: %v", err)
		}
	}

	if flushInterval == 0 {
		flushInterval = types.Duration(100 * time.Millisecond)
	}

	return &gohttputil.ReverseProxy{
		Director: func(outReq *http.Request) {
			u := outReq.URL
			if outReq.RequestURI != "/" {
				parsedURL, err := url.ParseRequestURI(outReq.RequestURI)
				if err == nil {
					u = parsedURL
				}
			}

			outReq.URL.Path = u.Path
			outReq.URL.RawPath = u.RawPath
			outReq.URL.RawQuery = u.RawQuery
			outReq.RequestURI = "" // Outgoing request should not have RequestURI

			outReq.Proto = "HTTP/1.1"
			outReq.ProtoMajor = 1
			outReq.ProtoMinor = 1

			// Do not pass client Host header unless optsetter PassHostHeader is set.
			if cfg.PassHostHeader != nil && !*cfg.PassHostHeader {
				outReq.Host = outReq.URL.Host
			}

			// It allows to proxy servers that uses authentication through URL (e.g. https://user:password@example.com)
			// In particular it allows to support nodes on Kaleido
			if u := outReq.URL.User; u != nil && outReq.Header.Get("Authorization") == "" {
				username := u.Username()
				password, _ := u.Password()
				outReq.Header.Set("Authorization", fmt.Sprintf("Basic %v", basicAuth(username, password)))
			}

			// Even if the websocket RFC says that headers should be case-insensitive,
			// some servers need Sec-WebSocket-Key to be case-sensitive.
			// https://tools.ietf.org/html/rfc6455#page-20
			outReq.Header["Sec-WebSocket-Key"] = outReq.Header["Sec-Websocket-Key"]
			delete(outReq.Header, "Sec-Websocket-Key")
		},
		Transport:      transport,
		FlushInterval:  time.Duration(flushInterval),
		ModifyResponse: respModifier,
		BufferPool:     pool,
		ErrorHandler: func(w http.ResponseWriter, request *http.Request, err error) {
			statusCode := http.StatusInternalServerError

			switch {
			case err == io.EOF:
				statusCode = http.StatusBadGateway
			case err == context.Canceled:
				statusCode = StatusClientClosedRequest
			default:
				if e, ok := err.(net.Error); ok {
					if e.Timeout() {
						statusCode = http.StatusGatewayTimeout
					} else {
						statusCode = http.StatusBadGateway
					}
				}
			}

			log.Debugf("'%d %s' caused by: %v", statusCode, statusText(statusCode), err)
			w.WriteHeader(statusCode)
			_, werr := w.Write([]byte(statusText(statusCode)))
			if werr != nil {
				log.Debugf("Error while writing status code", werr)
			}
		},
	}, nil
}

func statusText(statusCode int) string {
	if statusCode == StatusClientClosedRequest {
		return StatusClientClosedRequestText
	}
	return http.StatusText(statusCode)
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
