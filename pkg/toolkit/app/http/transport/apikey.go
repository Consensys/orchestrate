package transport

import (
	"net/http"

	authutils "github.com/consensys/orchestrate/pkg/toolkit/app/auth/utils"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
)

type XAPIKeyHeadersTransport struct {
	apiKey string
	T      http.RoundTripper
}

// NewXAPIKeyHeadersTransport creates a new transport to attach API-KEY as part of request headers
func NewXAPIKeyHeadersTransport(apiKey string) Middleware {
	return func(nxt http.RoundTripper) http.RoundTripper {
		return &XAPIKeyHeadersTransport{
			T:      nxt,
			apiKey: apiKey,
		}
	}
}

func (t *XAPIKeyHeadersTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	//FIXME CUSTOM HEADER double check this line, it may work to remove the GetAuthorizationHeader(req) == "" condition
	if authutils.GetAuthorizationHeader(req) == "" && t.apiKey != "" {
		authutils.AddAPIKeyHeaderValue(req, t.apiKey)
	}
	userInfo := multitenancy.UserInfoValue(req.Context())
	if userInfo != nil {
		if userInfo.TenantID != "" && userInfo.TenantID != multitenancy.WildcardTenant {
			authutils.AddTenantIDHeaderValue(req, userInfo.TenantID)
		}
		if userInfo.Username != "" {
			authutils.AddUsernameHeaderValue(req, userInfo.Username)
		}
	}

	return t.T.RoundTrip(req)
}
