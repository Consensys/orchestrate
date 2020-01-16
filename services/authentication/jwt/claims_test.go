package jwt

import (
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClaims(t *testing.T) {
	rawToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNYXBDbGFpbXMiOm51bGwsImh0dHA6Ly9vcmNoZXN0cmF0ZSI6eyJ0ZW5hbnRfaWQiOiJ0ZXN0LXRlbmFudCJ9fQ.Dwao4j6B_95PML5QL0TXY7Ys4rEDhDsrF4C5H6QTVKo"
	claims := &Claims{
		namespace: "http://orchestrate",
	}

	token, err := jwt.ParseWithClaims(rawToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("TestKey"), nil
	})

	require.NoError(t, err, "Should not fail in parsing token")
	require.True(t, token.Valid, "Token should be valid")
	assert.Equal(t, "test-tenant", token.Claims.(*Claims).Orchestrate.TenantID, "TenantID should be correct")
}
