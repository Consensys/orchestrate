package authentication

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope"

	"github.com/dgrijalva/jwt-go"
)

const (
	Component    = "authentication"
	TokenInfoKey = "token_info"
	TenantIDKey  = "tenant_id"
)

// Expose method to extract and verify the validity of an ID / Access Token (JWT format)
type Manager interface {
	// Extract the ID / Access Token from the envelop
	Extract(e *envelope.Envelope) (string, error)
	// Parse and verify the validity of the Token (ID or Access) and return a struct for a JWT (JSON Web Token)
	Verify(rawToken string) (*jwt.Token, error)
}
