package utils

import (
	authutils "github.com/consensys/orchestrate/pkg/toolkit/app/auth/utils"
)

func GetHeaders(apiKey, tenant, jwtToken string) map[string]string {
	headers := map[string]string{}
	if apiKey != "" {
		headers[authutils.APIKeyHeader] = apiKey
	}

	if tenant != "" {
		headers[authutils.TenantIDHeader] = tenant
	}

	if jwtToken != "" {
		headers[authutils.AuthorizationHeader] = jwtToken
	}

	return headers
}
