package token

import (
	"strings"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope"
)

const OauthToken = "oauth-token"

// Structure to define the parser of the Token and what have to be verify in the Token
type AuthToken struct {
	*jwt.Parser
}

func New() *AuthToken {
	return &AuthToken{
		&jwt.Parser{},
	}
}

// Extract the ID / Access Token from the envelop
func (a *AuthToken) Extract(e *envelope.Envelope) (string, error) {
	rawToken, ok := e.GetMetadataValue(OauthToken)
	if !ok {
		return "", errors.NotFoundError("JWT (ID Token or Access Token) was not found in the envelop").ExtendComponent(authentication.Component)
	}

	partsToken := strings.Split(rawToken, ".")
	if len(partsToken) != 3 {
		return "", errors.InvalidFormatError("JWT (ID Token or Access Token) is invalid token, token must have 3 partsToken").ExtendComponent(authentication.Component)
	}

	return rawToken, nil
}

// Parse and verify the validity of the Token (ID or Access) and return a struct for a JWT (JSON Web Token)
func (a *AuthToken) Verify(rawToken string) (*jwt.Token, error) {

	// Now parse the rawToken
	parsedToken, err := a.Parse(rawToken, authentication.ValidatedKey)

	// Check if there was an error in parsing...
	if err != nil {
		return parsedToken, err
	}

	// Check if the parsed rawToken is valid...
	if !parsedToken.Valid {
		return parsedToken, errors.InvalidCryptographicSignatureError("token is invalid").ExtendComponent(authentication.Component)
	}

	log.Tracef("ID / Access Token is validate")

	return parsedToken, nil
}
