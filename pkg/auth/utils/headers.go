package utils

import (
	"net/http"
)

const (
	AuthorizationHeader = "Authorization"
	APIKeyHeader        = "X-API-Key"
)

func AddXAPIKeyHeaderValue(req *http.Request, apiKey string) {
	req.Header.Add(APIKeyHeader, apiKey)
}

func AddAuthorizationHeader(req *http.Request) {
	authorization := AuthorizationFromContext(req.Context())
	if authorization != "" {
		req.Header.Add(AuthorizationHeader, authorization)
	}
}

func GetAuthorizationHeader(req *http.Request) string {
	return req.Header.Get(AuthorizationHeader)
}
