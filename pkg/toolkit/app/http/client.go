package http

import (
	"crypto/tls"
	"net"
	"net/http"

	"github.com/consensys/orchestrate/pkg/toolkit/app/http/transport"
)

// NewClient creates an HTTP client
func NewClient(cfg *Config) *http.Client {
	dialer := &net.Dialer{
		Timeout:   cfg.Timeout,
		KeepAlive: cfg.KeepAlive,
	}

	/** Execution flow
	1. Attach Authentication Headers, if they are part of context
	2. Attach X-API-KEY header, only if Authentication was not set before
	2. Attach Authorization header, only if Authentication was not set before
	3. Retry on 429 responses
	*/
	middlewares := []transport.Middleware{}
	if cfg.Authorization != "" {
		middlewares = append(middlewares, transport.NewAuthHeadersTransport(cfg.Authorization))
	}

	if cfg.AuthHeaderForward {
		middlewares = append(middlewares, transport.NewContextAuthHeadersTransport())
	}

	if cfg.XAPIKey != "" {
		middlewares = append(middlewares, transport.NewXAPIKeyHeadersTransport(cfg.XAPIKey))
	}

	middlewares = append(middlewares, transport.NewRetry429Transport())

	t := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		MaxIdleConnsPerHost:   cfg.MaxIdleConnsPerHost,
		IdleConnTimeout:       cfg.IdleConnTimeout,
		TLSHandshakeTimeout:   cfg.TLSHandshakeTimeout,
		ExpectContinueTimeout: cfg.ExpectContinueTimeout,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: cfg.InsecureSkipVerify,
		},
	}

	if cfg.ClientCert != nil {
		t.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: cfg.InsecureSkipVerify,
			Certificates:       []tls.Certificate{*cfg.ClientCert},
			GetClientCertificate: func(info *tls.CertificateRequestInfo) (*tls.Certificate, error) {
				return cfg.ClientCert, nil
			},
		}
	}

	return &http.Client{
		Transport: NewTransport(t, middlewares...),
	}
}
