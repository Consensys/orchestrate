package utils

import (
	"net/http"
)

const (
	AuthorizationHeader = "Authorization"
	APIKeyHeader        = "X-API-Key"
)

func AddXAPIKeyHeader(req *http.Request) {
	apiKey := APIKeyFromContext(req.Context())
	if apiKey != "" {
		req.Header.Add(APIKeyHeader, apiKey)
	}
}

func AddAuthorizationHeader(req *http.Request) {
	authorization := AuthorizationFromContext(req.Context())
	if authorization != "" {
		req.Header.Add(AuthorizationHeader, authorization)
	}
}
