package authentication

import (
	"github.com/dgrijalva/jwt-go"
)

const (
	Component    = "authentication"
	TokenRawKey  = "token_raw"
	TokenInfoKey = "token_info"
	TenantIDKey  = "tenant_id"
)

// Expose method to extract and verify the validity of an ID / Access Token (JWT format)
type Manager interface {
	// Parse and verify the validity of the Token (ID or Access) and return a struct for a JWT (JSON Web Token)
	Verify(rawToken string) (*jwt.Token, error)
}
