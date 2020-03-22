package multitenancy

import (
	"net/http"
)

// The TenantID Header have to be used only between tx-listener and envelope-store
const TenantIDHeader = "X-Tenant-ID"

func AddTenantIDHeader(req *http.Request) {
	tenantID := TenantIDFromContext(req.Context())
	if tenantID != "" {
		req.Header.Add(TenantIDHeader, tenantID)
	}
}
