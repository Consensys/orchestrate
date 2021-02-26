package proxy

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
	"time"

	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"github.com/containous/traefik/v2/pkg/log"
	traefiktls "github.com/containous/traefik/v2/pkg/tls"
	"golang.org/x/net/http2"
)

type h2cTransportWrapper struct {
	*http2.Transport
}

func (t *h2cTransportWrapper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "http"
	return t.Transport.RoundTrip(req)
}

// NewTransport creates an http.RoundTripper configured with the Transport configuration settings.
// For the settings that can't be configured in Traefik it uses the default http.Transport settings.
// An exception to this is the MaxIdleConns setting as we only provide the option MaxIdleConnsPerHost
// in Traefik at this point in time. Setting this value to the default of 100 could lead to confusing
// behavior and backwards compatibility issues.
func NewTransport(cfg *traefikstatic.ServersTransport) (http.RoundTripper, error) {
	if cfg == nil {
		return http.DefaultTransport, nil
	}

	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	if cfg.ForwardingTimeouts != nil {
		dialer.Timeout = time.Duration(cfg.ForwardingTimeouts.DialTimeout)
	}

	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		MaxIdleConnsPerHost:   cfg.MaxIdleConnsPerHost,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	transport.RegisterProtocol("h2c", &h2cTransportWrapper{
		Transport: &http2.Transport{
			DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(netw, addr)
			},
			AllowHTTP: true,
		},
	})

	if cfg.ForwardingTimeouts != nil {
		transport.ResponseHeaderTimeout = time.Duration(cfg.ForwardingTimeouts.ResponseHeaderTimeout)
		transport.IdleConnTimeout = time.Duration(cfg.ForwardingTimeouts.IdleConnTimeout)
	}

	if cfg.InsecureSkipVerify {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	if len(cfg.RootCAs) > 0 {
		transport.TLSClientConfig = &tls.Config{
			RootCAs: createRootCACertPool(cfg.RootCAs),
		}
	}

	err := http2.ConfigureTransport(transport)
	if err != nil {
		return nil, err
	}

	return transport, nil
}

func createRootCACertPool(rootCAs []traefiktls.FileOrContent) *x509.CertPool {
	roots := x509.NewCertPool()

	for _, cert := range rootCAs {
		certContent, err := cert.Read()
		if err != nil {
			log.WithoutContext().Error("Error while read RootCAs", err)
			continue
		}
		roots.AppendCertsFromPEM(certContent)
	}

	return roots
}
