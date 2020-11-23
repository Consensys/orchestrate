package header

import (
	"net/textproto"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/auth/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
)

var (
	authorizationHeader = textproto.CanonicalMIMEHeaderKey(authutils.AuthorizationHeader)
	apiKeyHeader        = textproto.CanonicalMIMEHeaderKey(authutils.APIKeyHeader)
	tenantIDHeader      = textproto.CanonicalMIMEHeaderKey(multitenancy.TenantIDHeader)
)

// AuthCredMatcher Verifies that the header is part of the accepted headers
func AuthCredMatcher(key string) (string, bool) {
	key = textproto.CanonicalMIMEHeaderKey(key)
	if key == authorizationHeader || key == apiKeyHeader || key == tenantIDHeader {
		return key, true
	}
	return runtime.DefaultHeaderMatcher(key)
}
