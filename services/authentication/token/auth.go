package token

import (
	"context"
	"strings"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope"
)

const (
	OauthToken = "oauth-token"
	HeaderKey  = "Bearer"
)

// Structure to define the parser of the Token and what have to be verify in the Token
type AuthToken struct {
	*jwt.Parser
}

func New() *AuthToken {
	return &AuthToken{
		&jwt.Parser{},
	}
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

	log.Tracef("ID/Access Token is valid")

	return parsedToken, nil
}

// Extract the ID / Access Token from the envelop
func Extract(e *envelope.Envelope) (string, error) {
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

// Extract JWT from Envelop to inject it into client gRPC
func GetGRPCOptionJWTokenFromEnvelope(txctx *engine.TxContext) grpc.CallOption {
	rawToken, err := Extract(txctx.Envelope)
	if err != nil {
		rawToken = ""
	}

	return NewJWTokenGRPCOption(rawToken)
}

func GetGRPCOptionJWTokenFromContext(ctx context.Context) grpc.CallOption {

	rawToken, ok := ctx.Value(authentication.TokenRawKey).(string)
	if !ok {
		rawToken = ""
	}

	return NewJWTokenGRPCOption(rawToken)
}

func NewJWTokenGRPCOption(rawToken string) grpc.CallOption {

	perRPCCredentials := NewJWTAccessFromEnvelope(rawToken)
	return grpc.PerRPCCredentials(perRPCCredentials)
}

type jwtAccess struct {
	jwToken string
}

func (j jwtAccess) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	var headerJWT map[string]string

	if j.jwToken != "" {
		headerJWT = map[string]string{
			"authorization": HeaderKey + " " + j.jwToken,
		}
	} else {
		headerJWT = map[string]string{}
	}

	return headerJWT, nil
}

func (j jwtAccess) RequireTransportSecurity() bool {
	return false
}

// NewJWTAccessFromKey creates PerRPCCredentials from the given token.
func NewJWTAccessFromEnvelope(rawToken string) credentials.PerRPCCredentials {
	return jwtAccess{rawToken}
}
