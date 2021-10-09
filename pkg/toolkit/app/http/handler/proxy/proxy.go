package proxy

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	gohttputil "net/http/httputil"
	"net/url"
	"time"

	"github.com/consensys/orchestrate/pkg/toolkit/app/http/config/dynamic"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/httputil"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/middleware/loadbalancer"
	"github.com/oxtoacart/bpool"
	traefiktypes "github.com/traefik/paerser/types"
	traefikstatic "github.com/traefik/traefik/v2/pkg/config/static"
	"github.com/traefik/traefik/v2/pkg/log"
)

// StatusClientClosedRequest non-standard HTTP status code for client disconnection
const StatusClientClosedRequest = 499

// StatusClientClosedRequestText non-standard HTTP status for client disconnection
const StatusClientClosedRequestText = "Client Closed Connection"

type Builder struct {
	transport http.RoundTripper
	bpool     gohttputil.BufferPool

	lb *loadbalancer.Builder
}

func NewBuilder(transportCfg *traefikstatic.ServersTransport, pool gohttputil.BufferPool) (*Builder, error) {
	t, err := NewTransport(transportCfg)
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
	var flushInterval traefiktypes.Duration
	if cfg.ResponseForwarding != nil && cfg.ResponseForwarding.FlushInterval != "" {
		err := flushInterval.Set(cfg.ResponseForwarding.FlushInterval)
		if err != nil {
			return nil, fmt.Errorf("error creating flush interval: %v", err)
		}
	}

	if flushInterval == 0 {
		flushInterval = traefiktypes.Duration(100 * time.Millisecond)
	}

	return &gohttputil.ReverseProxy{
		Director: func(outReq *http.Request) {
			u := outReq.URL

			if outReq.RequestURI != "/" {
				// Add requestURI in the path
				// In case the downstream backend is domain/path (i.e. domain.com/tessera), and the request is done to proxy/proxyPath/(.+) (i.e. proxy.com/proxyPath/storeraw)
				// Then backend call will be domain.com/tessera/storeraw
				// if u.EscapedPath() != outReq.RequestURI {
				// 	outReq.RequestURI = u.EscapedPath() + outReq.RequestURI
				// }

				parsedURL, err := url.ParseRequestURI(u.EscapedPath() + outReq.RequestURI)
				if err == nil {
					u = parsedURL
				}
			}

			outReq.URL.Path = u.Path
			outReq.URL.RawPath = u.RawPath
			outReq.URL.RawQuery = u.RawQuery
			outReq.RequestURI = "" // Outgoing request should not have RequestURI

			// See More https://github.com/golang/net/blob/master/http2/transport.go#L554
			if outReq.Body != nil {
				body, _ := ioutil.ReadAll(outReq.Body)
				outReq.Body = ioutil.NopCloser(bytes.NewBuffer(body))
				outReq.GetBody = func() (io.ReadCloser, error) {
					return ioutil.NopCloser(bytes.NewBuffer(body)), nil
				}
			} else {
				outReq.GetBody = func() (io.ReadCloser, error) {
					return nil, nil
				}
			}

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
		Transport:     transport,
		FlushInterval: time.Duration(flushInterval),
		// ModifyResponse: respModifier,
		ModifyResponse: func(rw *http.Response) error {
			rw.Header.Set("X-Backend-Server", rw.Request.URL.String())
			if rw.StatusCode >= 300 {
				body, _ := ioutil.ReadAll(rw.Body)
				rw.Body = ioutil.NopCloser(bytes.NewBuffer(body))
				log.FromContext(rw.Request.Context()).
					Debugf("'%d %s' caused by: %q", rw.StatusCode, statusText(rw.StatusCode), string(body))
			}

			if respModifier == nil {
				return nil
			}

			return respModifier(rw)
		},
		BufferPool: pool,
		ErrorHandler: func(w http.ResponseWriter, req *http.Request, err error) {
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
						statusCode = http.StatusServiceUnavailable
					}
				}
			}

			log.FromContext(req.Context()).Debugf("'%d %s' caused by: %v", statusCode, statusText(statusCode), err)

			w.Header().Set("X-Backend-Server", req.URL.String())
			w.WriteHeader(statusCode)
			_, werr := w.Write([]byte(statusText(statusCode)))
			if werr != nil {
				log.FromContext(req.Context()).Debugf("Error while writing status code", werr)
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
