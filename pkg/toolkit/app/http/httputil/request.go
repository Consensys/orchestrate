package httputil

import (
	"crypto/tls"
	"html"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"

	traefiktls "github.com/traefik/traefik/v2/pkg/tls"
)

func GetMethod(r *http.Request) string {
	if !utf8.ValidString(r.Method) {
		return "NON_UTF8_HTTP_METHOD"
	}
	return r.Method
}

func GetTLSVersion(req *http.Request) string {
	if req.TLS == nil {
		panic("request TLS config is not set")
	}

	switch req.TLS.Version {
	case tls.VersionTLS10:
		return "1.0"
	case tls.VersionTLS11:
		return "1.1"
	case tls.VersionTLS12:
		return "1.2"
	case tls.VersionTLS13:
		return "1.3"
	default:
		return "UNKNOWN_TLS_VERSION"
	}
}

func GetTLSCipher(req *http.Request) string {
	if req.TLS == nil {
		panic("request TLS config is not set")
	}

	if version, ok := traefiktls.CipherSuitesReversed[req.TLS.CipherSuite]; ok {
		return version
	}

	return "UNKNOWN_TLS_CIPHER"
}

func GetProtocol(req *http.Request) string {
	switch {
	case IsWebsocketRequest(req):
		return "websocket"
	case IsSSERequest(req):
		return "sse"
	default:
		return "http"
	}
}

// isWebsocketRequest determines if the specified HTTP request is a websocket handshake request.
func IsWebsocketRequest(req *http.Request) bool {
	return ContainsHeader(req, "Connection", "upgrade") && ContainsHeader(req, "Upgrade", "websocket")
}

// isSSERequest determines if the specified HTTP request is a request for an event subscription.
func IsSSERequest(req *http.Request) bool {
	return ContainsHeader(req, "Accept", "text/event-stream")
}

func ContainsHeader(req *http.Request, name, value string) bool {
	items := strings.Split(req.Header.Get(name), ",")
	for _, item := range items {
		if strings.EqualFold(value, strings.TrimSpace(item)) {
			return true
		}
	}
	return false
}

func ToFilters(values url.Values) map[string]string {
	filters := make(map[string]string)
	for key := range values {
		k := html.EscapeString(key)
		v := html.EscapeString(values.Get(key))
		if k != "" && v != "" {
			filters[k] = v
		}
	}
	return filters
}
