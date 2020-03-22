// +build unit

package header

import (
	"net/textproto"
	"testing"

	"github.com/stretchr/testify/assert"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
)

func TestAuthCredMatcher(t *testing.T) {
	headerAuthorization, _ := AuthCredMatcher(authutils.AuthorizationHeader)
	assert.Equal(t, textproto.CanonicalMIMEHeaderKey(authutils.AuthorizationHeader), headerAuthorization)

	headerTenantID, _ := AuthCredMatcher(multitenancy.TenantIDHeader)
	assert.Equal(t, textproto.CanonicalMIMEHeaderKey(multitenancy.TenantIDHeader), headerTenantID)

	headerAPIKey, _ := AuthCredMatcher(authutils.APIKeyHeader)
	assert.Equal(t, textproto.CanonicalMIMEHeaderKey(authutils.APIKeyHeader), headerAPIKey)

	headerInvalid, ok := AuthCredMatcher("InvalidHeader")
	assert.Equal(t, "", headerInvalid)
	assert.False(t, ok)
}
