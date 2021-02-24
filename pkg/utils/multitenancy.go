package utils

import (
	"github.com/ConsenSys/orchestrate/pkg/multitenancy"
)

const (
	AuthorizationMetadata = "Authorization"
	TenantIDMetadata      = "X-Tenant-ID"
)

func AllowedTenants(tenantID string) []string {
	// If no tenant then we use default tenant
	if tenantID == multitenancy.DefaultTenant {
		return []string{multitenancy.DefaultTenant}
	}

	// Otherwise the account can belong either to the default tenant or to the specified one
	return []string{multitenancy.DefaultTenant, tenantID}
}
