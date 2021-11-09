package utils

import (
	authutils "github.com/consensys/orchestrate/pkg/toolkit/app/auth/utils"
	"github.com/consensys/orchestrate/pkg/utils"
)

func GetHeaders(apiKey, tenant, jwtToken string) map[string]string {
	headers := map[string]string{}
	if apiKey != "" {
		headers[authutils.APIKeyHeader] = apiKey
	}

	if tenant != "" {
		headers[utils.TenantIDMetadata] = tenant
	}

	if jwtToken != "" {
		headers[utils.AuthorizationMetadata] = jwtToken
	}

	return headers
}
