package http

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"github.com/containous/traefik/v2/pkg/log"
	"github.com/containous/traefik/v2/pkg/middlewares/forwardedheaders"
	"github.com/hashicorp/go-multierror"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/handler/switcher"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/router"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/tcp"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/tcp/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const (
	DefaultHTTPAppEntryPoint = "app"
	DefaultMetricsEntryPoint = "metrics"
)

type EntryPoints struct {
	eps    map[string]*EntryPoint
	router router.Builder
	errors chan error
}

func NewEntryPoints(
	epConfigs traefikstatic.EntryPoints,
	rt router.Builder,
	reg metrics.TPCMetrics,
) *EntryPoints {
	s := &EntryPoints{
		eps:    make(map[string]*EntryPoint),
		router: rt,
		errors: make(chan error, len(epConfigs)),
	}

	for epName, epConfig := range epConfigs {
		var middlewares []func(h http.Handler) http.Handler
		if epConfig.ForwardedHeaders != nil {
			mid := func(h http.Handler) http.Handler {
				m, _ := forwardedheaders.NewXForwarded(
					epConfig.ForwardedHeaders.Insecure,
					epConfig.ForwardedHeaders.TrustedIPs,
					h,
				)
				return m
			}
			middlewares = append(middlewares, mid)
		}

		httpServer := newSwitchableServer(epConfig.Transport.RespondingTimeouts, true, middlewares...)
		httpsServer := newSwitchableServer(epConfig.Transport.RespondingTimeouts, false, middlewares...)

		s.eps[epName] = NewEntryPoint(epName, epConfig, httpServer, httpsServer, reg)
	}

	return s
}

func (eps *EntryPoints) Addresses() map[string]string {
	addrs := make(map[string]string)
	for epName, ep := range eps.eps {
		addrs[epName] = ep.Addr()
	}
	return addrs
}

func (eps *EntryPoints) ListenAndServe(ctx context.Context) chan error {
	errors := make(chan error, len(eps.eps))
	go func() {
		wg := &sync.WaitGroup{}
		wg.Add(len(eps.eps))
		for _, ep := range eps.eps {
			go func(ep *EntryPoint) {
				if err := ep.ListenAndServe(ctx); err != http.ErrServerClosed {
					errors <- err
				}
				wg.Done()
			}(ep)
		}
		wg.Wait()
		close(errors)
	}()

	return errors
}

func (eps *EntryPoints) Errors() <-chan error {
	return eps.errors
}

func (eps *EntryPoints) Switch(ctx context.Context, conf interface{}) error {
	var entryPointNames []string
	for epName := range eps.eps {
		entryPointNames = append(entryPointNames, epName)
	}

	rt, err := eps.router.Build(ctx, entryPointNames, conf)
	if err != nil {
		log.FromContext(ctx).WithError(err).Errorf("error building router")
		return err
	}

	eps.switchRouter(ctx, rt)

	return nil
}

func (eps *EntryPoints) switchRouter(ctx context.Context, routers map[string]*router.Router) {
	for epName, ep := range eps.eps {
		logger := log.FromContext(ctx).WithField("entrypoint", epName)
		rt, ok := routers[epName]
		if ok {
			err := ep.Switch(rt)
			if err != nil {
				logger.WithError(err).Errorf("error switching tcp router")
			} else {
				logger.Infof("switched tcp router")
			}
		}
	}
}

func (eps *EntryPoints) Shutdown(ctx context.Context) error {
	gr := &multierror.Group{}
	for epName, ep := range eps.eps {
		epName, ep := epName, ep
		gr.Go(func() error { return tcp.Shutdown(log.With(ctx, log.Str("entrypoint", epName)), ep) })
	}

	return gr.Wait().ErrorOrNil()
}

func (eps *EntryPoints) Close() error {
	gr := &multierror.Group{}
	for _, ep := range eps.eps {
		ep := ep
		gr.Go(func() error { return tcp.Close(ep) })
	}
	return gr.Wait().ErrorOrNil()
}

type EntryPoint struct {
	cfg      *traefikstatic.EntryPoint
	tcp      *tcp.EntryPoint
	switcher *switchTCPHandler
}

func NewEntryPoint(name string, ep *traefikstatic.EntryPoint, httpSrv, httpsSrv *switchableServer, reg metrics.TPCMetrics) *EntryPoint {
	tcpSwitcher := &switchTCPHandler{
		switcher:       tcp.NewSwitcher(),
		http:           httpSrv,
		https:          httpsSrv,
		httpForwarder:  tcp.NewForwarder(nil),
		httpsForwarder: tcp.NewForwarder(nil),
	}

	return &EntryPoint{
		cfg: ep,
		tcp: tcp.NewEntryPoint(
			name,
			ep, tcpSwitcher,
			reg,
		),
		switcher: tcpSwitcher,
	}
}

func (ep *EntryPoint) Addr() string {
	return ep.tcp.Addr()
}

func (ep *EntryPoint) ListenAndServe(ctx context.Context) error {
	go func() {
		_ = ep.switcher.ListenAndServe()
	}()

	return ep.tcp.ListenAndServe(ctx)
}

func (ep *EntryPoint) Switch(rt *router.Router) error {
	ep.switcher.Switch(rt)
	return nil
}

func (ep *EntryPoint) Shutdown(ctx context.Context) error {
	return tcp.Shutdown(ctx, ep.tcp)
}

func (ep *EntryPoint) Close() error {
	return tcp.Close(ep.tcp)
}

type switchTCPHandler struct {
	switcher *tcp.Switcher

	http  *switchableServer
	https *switchableServer

	httpForwarder  *tcp.Forwarder
	httpsForwarder *tcp.Forwarder
}

func (s *switchTCPHandler) ServeTCP(conn tcp.WriteCloser) {
	s.switcher.ServeTCP(conn)
}

func (s *switchTCPHandler) ListenAndServe() error {
	// We can ignore next errors since any net.Error are catched
	// at the tcp.EntryPoint level
	utils.InParallel(
		func() { _ = s.http.serve(s.httpForwarder) },
		func() { _ = s.https.serve(s.httpsForwarder) },
	)

	return http.ErrServerClosed
}

func (s *switchTCPHandler) Switch(conf *router.Router) {
	rt := &tcp.Router{}

	// Set router TLS configuration
	rt.TLSConfig(conf.TLSConfig)
	for sni, tlsConfig := range conf.HostTLSConfigs {
		rt.AddRouteHTTPTLS(sni, tlsConfig)
	}

	// Set forwarders for HTTP & HTTPS server
	rt.TCPForwarder(s.httpForwarder)
	rt.TLSForwarder(s.httpsForwarder)

	// Switch Handlers on HTTP & HTTPS servers
	s.http.switchHandler(conf.HTTP)
	s.https.switchHandler(conf.HTTPS)

	// Switch router
	s.switcher.Switch(rt)
}

func (s *switchTCPHandler) Shutdown(ctx context.Context) error {
	gr := &multierror.Group{}
	gr.Go(func() error { return tcp.Shutdown(ctx, s.http) })
	gr.Go(func() error { return tcp.Shutdown(ctx, s.https) })
	return gr.Wait().ErrorOrNil()
}

func (s *switchTCPHandler) Close() error {
	gr := &multierror.Group{}
	gr.Go(func() error { return tcp.Close(s.http) })
	gr.Go(func() error { return tcp.Close(s.https) })
	gr.Go(func() error { return tcp.Close(s.httpForwarder) })
	gr.Go(func() error { return tcp.Close(s.httpsForwarder) })
	return gr.Wait().ErrorOrNil()
}

type switchableServer struct {
	server   *http.Server
	switcher *switcher.Switcher
}

func newSwitchableServer(
	timeouts *traefikstatic.RespondingTimeouts,
	withH2c bool,
	middlewares ...func(http.Handler) http.Handler,
) *switchableServer {
	swtchr := switcher.New()
	var handler http.Handler = swtchr
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}

	if withH2c {
		handler = h2c.NewHandler(handler, &http2.Server{})
	}

	server := &http.Server{
		Handler:      handler,
		ReadTimeout:  time.Duration(timeouts.ReadTimeout),
		WriteTimeout: time.Duration(timeouts.WriteTimeout),
		IdleTimeout:  time.Duration(timeouts.IdleTimeout),
	}

	return &switchableServer{
		server:   server,
		switcher: swtchr,
	}
}

func (s *switchableServer) serve(l net.Listener) error {
	return s.server.Serve(l)
}

func (s *switchableServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *switchableServer) Close() error {
	return s.server.Close()
}

func (s *switchableServer) switchHandler(h http.Handler) {
	s.switcher.Switch(h)
}
