package accesslog

import (
	"net/http"
)

const (
	// StartUTC is the map key used for the time at which request processing started.
	StartUTC = "process.start"
	// StartLocal is the map key used for the local time at which request processing started.
	StartLocal = "process.start_local" // not a standard
	// Duration is the map key used for the total time taken by processing the response, including the origin server's time but
	// not the log writing time.
	Duration = "process.uptime"

	// RouterName is the map key used for the name of the Traefik router.
	RouterName = "service.router_name" // not a standard
	// ServiceName is the map key used for the name of the Traefik backend.
	ServiceName = "service.name"
	// ServiceURL is the map key used for the URL of the Traefik backend.
	ServiceURL = "service.url"
	// ServiceAddr is the map key used for the IP:port of the Traefik backend (extracted from BackendURL)
	ServiceAddr = "service.addr"

	// ClientAddr is the map key used for the remote address in its original form (usually IP:port).
	ClientAddr = "client.addr"
	// ClientHost is the map key used for the remote IP address from which the client request was received.
	ClientHost = "client.host"
	// ClientPort is the map key used for the remote TCP port from which the client request was received.
	ClientPort = "client.port"
	// ClientUsername is the map key used for the username provided in the URL, if present.
	ClientUsername = "client.user.name"
	// RequestAddr is the map key used for the HTTP Host header (usually IP:port). This is treated as not a header by the Go API.
	RequestAddr = "url.full"
	// RequestHost is the map key used for the HTTP Host server name (not including port).
	RequestHost = "url.domain"
	// RequestPort is the map key used for the TCP port from the HTTP Host.
	RequestPort = "url.port"
	// RequestMethod is the map key used for the HTTP method.
	RequestMethod = "http.request.method"
	// RequestPath is the map key used for the HTTP request URI, not including the scheme, host or port.
	RequestPath = "url.path"
	// RequestProtocol is the map key used for the version of HTTP requested.
	RequestProtocol = "http.request.protocol" // not a standard
	// RequestScheme is the map key used for the HTTP request scheme.
	RequestScheme = "url.scheme"
	// RequestContentSize is the map key used for the number of bytes in the request entity (a.k.a. body) sent by the client.
	RequestContentSize = "http.request.body.bytes"
	// RequestRefererHeader is the Referer header in the request
	RequestRefererHeader = "request_Referer"
	// RequestUserAgentHeader is the User-Agent header in the request
	RequestUserAgentHeader = "request_User-Agent"
	// OriginDuration is the map key used for the time taken by the origin server ('upstream') to return its response.
	OriginDuration = "source.duration" // not a standard
	// OriginContentSize is the map key used for the content length specified by the origin server, or 0 if unspecified.
	OriginContentSize = "source.bytes"
	// OriginStatus is the map key used for the HTTP status code returned by the origin server.
	// If the request was handled by this Traefik instance (e.g. with a redirect), then this value will be absent.
	OriginStatus = "source.http.response.status_code" // not a standard
	// DownstreamStatus is the map key used for the HTTP status code returned to the client.
	DownstreamStatus = "http.response.status_code"
	// DownstreamContentSize is the map key used for the number of bytes in the response entity returned to the client.
	// This is in addition to the "Content-Length" header, which may be present in the origin response.
	DownstreamContentSize = "http.response.bytes"
	// RequestCount is the map key used for the number of requests received since the Traefik instance started.
	RequestCount = "http.request.count" // not a standard
	// GzipRatio is the map key used for the response body compression ratio achieved.
	GzipRatio = "http.response.gzip_ratio" // not a standard
	// Overhead is the map key used for the processing time overhead caused by Traefik.
	Overhead = "process.overhead" // not a standard
	// RetryAttempts is the map key used for the amount of attempts the request was retried.
	RetryAttempts = "http.request.retry_attempts" // not a standard
)

// These are written out in the default case when no config is provided to specify keys of interest.
var defaultCoreKeys = [...]string{
	StartUTC,
	Duration,
	RouterName,
	ServiceName,
	ServiceURL,
	ClientHost,
	ClientPort,
	ClientUsername,
	RequestHost,
	RequestPort,
	RequestMethod,
	RequestPath,
	RequestProtocol,
	RequestScheme,
	RequestContentSize,
	OriginDuration,
	OriginContentSize,
	OriginStatus,
	DownstreamStatus,
	DownstreamContentSize,
	RequestCount,
}

// This contains the set of all keys, i.e. all the default keys plus all non-default keys.
var allCoreKeys = make(map[string]struct{})

func init() {
	for _, k := range defaultCoreKeys {
		allCoreKeys[k] = struct{}{}
	}
	allCoreKeys[ServiceAddr] = struct{}{}
	allCoreKeys[ClientAddr] = struct{}{}
	allCoreKeys[RequestAddr] = struct{}{}
	allCoreKeys[GzipRatio] = struct{}{}
	allCoreKeys[StartLocal] = struct{}{}
	allCoreKeys[Overhead] = struct{}{}
	allCoreKeys[RetryAttempts] = struct{}{}
}

// CoreLogData holds the fields computed from the request/response.
type CoreLogData map[string]interface{}

// LogData is the data captured by the middleware so that it can be logged.
type LogData struct {
	Core               CoreLogData
	Request            request
	OriginResponse     http.Header
	DownstreamResponse downstreamResponse
	Message            string
}

type downstreamResponse struct {
	headers http.Header
	status  int
	size    int64
}

type request struct {
	headers http.Header
	// Request body size
	size int64
}
