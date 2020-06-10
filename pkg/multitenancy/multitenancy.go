package multitenancy

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

const Wildcard = "*"

func IsAllowed(tenantID string, tenants []string) bool {
	for _, tenant := range tenants {
		if tenantID == tenant || tenant == Wildcard {
			return true
		}
	}
	return false
}

func AllowedTenants(jwtTenant, tenantID string) []string {
	switch {
	case jwtTenant == "":
		return []string{}
	case jwtTenant == Wildcard && tenantID == "":
		return []string{Wildcard}
	case jwtTenant == Wildcard:
		return []string{tenantID}
	case tenantID == DefaultTenant:
		return []string{DefaultTenant}
	case tenantID == "":
		return []string{DefaultTenant, jwtTenant}
	case jwtTenant != tenantID:
		return []string{}
	default:
		return []string{jwtTenant}
	}
}

func TenantID(jwtTenant, tenantID string) (string, error) {
	switch {
	case jwtTenant == "":
		return "", errors.PermissionDeniedError("empty tenant in Access Token")
	case jwtTenant == Wildcard && tenantID == "":
		return DefaultTenant, nil
	case jwtTenant == Wildcard:
		return tenantID, nil
	case tenantID == DefaultTenant:
		return DefaultTenant, nil
	case tenantID == "":
		return jwtTenant, nil
	case jwtTenant != tenantID:
		return "", errors.PermissionDeniedError("access to tenant %q forbidden", tenantID)
	default:
		return jwtTenant, nil
	}
}
