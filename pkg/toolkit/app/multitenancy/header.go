package multitenancy

import (
	"net/http"
)

const TenantIDHeader = "X-Tenant-ID"

func AddTenantIDHeader(req *http.Request) {
	tenantID, ok := TenantIDValue(req.Context())
	if ok {
		req.Header.Add(TenantIDHeader, tenantID)
	}
}
