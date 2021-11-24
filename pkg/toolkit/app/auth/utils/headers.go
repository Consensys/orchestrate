package utils

import (
	"net/http"
)

const (
	AuthorizationHeader = "Authorization"
	APIKeyHeader        = "X-API-Key"
	TenantIDHeader      = "X-Tenant-ID"
	UsernameHeader      = "X-Username"
)

func GetAPIKeyHeaderValue(req *http.Request) string {
	return req.Header.Get(APIKeyHeader)
}

func AddAPIKeyHeaderValue(req *http.Request, value string) {
	req.Header.Add(APIKeyHeader, value)
}

func DeleteAPIKeyHeaderValue(req *http.Request) {
	req.Header.Del(APIKeyHeader)
}

func GetAuthorizationHeader(req *http.Request) string {
	return req.Header.Get(AuthorizationHeader)
}

func AddAuthorizationHeaderValue(req *http.Request, value string) {
	req.Header.Add(AuthorizationHeader, value)
}

func DeleteAuthorizationHeaderValue(req *http.Request) {
	req.Header.Del(AuthorizationHeader)
}

func GetTenantIDHeaderValue(req *http.Request) string {
	return req.Header.Get(TenantIDHeader)
}

func AddTenantIDHeaderValue(req *http.Request, value string) {
	req.Header.Add(TenantIDHeader, value)
}

func GetUsernameHeaderValue(req *http.Request) string {
	return req.Header.Get(UsernameHeader)
}

func AddUsernameHeaderValue(req *http.Request, value string) {
	req.Header.Add(UsernameHeader, value)
}
