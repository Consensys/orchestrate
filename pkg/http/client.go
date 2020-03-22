package http

import (
	"net"
	"net/http"
	"time"
)

// NewClient creates an HTTP client
func NewClient() *http.Client {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}

	return &http.Client{
		Transport: NewTransport(&http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           dialer.DialContext,
			MaxIdleConnsPerHost:   200,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}),
	}
}
