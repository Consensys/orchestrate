package http

import (
	"net"
	"net/http"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/transport"
)

// NewClient creates an HTTP client
func NewClient(cfg *Config) *http.Client {
	dialer := &net.Dialer{
		Timeout:   cfg.Timeout,
		KeepAlive: cfg.KeepAlive,
	}

	/** Execution flow
	1. (only multi-tenancy) Attach Authentication Headers, if they are part of context
	2. (only multi-tenancy) Attach API-KEY headers, only if Authentication was not set before
	3. Retry on 429 responses
	*/
	middlewares := []transport.Middleware{}
	if cfg.MultiTenancy {
		middlewares = append(middlewares,
			transport.NewAuthHeadersTransport(),
			transport.NewAPIKeyHeadersTransport(cfg.APIKey),
		)
	}

	middlewares = append(middlewares, transport.NewRetry429Transport())

	return &http.Client{
		Transport: NewTransport(&http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           dialer.DialContext,
			MaxIdleConnsPerHost:   cfg.MaxIdleConnsPerHost,
			IdleConnTimeout:       cfg.IdleConnTimeout,
			TLSHandshakeTimeout:   cfg.TLSHandshakeTimeout,
			ExpectContinueTimeout: cfg.ExpectContinueTimeout,
		}, middlewares...),
	}
}
