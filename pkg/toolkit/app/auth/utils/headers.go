package utils

import (
	"net/http"
)

const (
	AuthorizationHeader = "Authorization"
	APIKeyHeader        = "X-API-Key"
	TenantIDHeader      = "X-Tenant-ID"
)

func AddXAPIKeyHeaderValue(req *http.Request, apiKey string) {
	req.Header.Add(APIKeyHeader, apiKey)
}

func AddAuthorizationHeaderValue(req *http.Request, authorization string) {
	req.Header.Add(AuthorizationHeader, authorization)
}

func GetAuthorizationHeader(req *http.Request) string {
	return req.Header.Get(AuthorizationHeader)
}
