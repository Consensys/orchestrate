package jwt

import (
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	"github.com/consensys/orchestrate/pkg/errors"
	authutils "github.com/consensys/orchestrate/pkg/toolkit/app/auth/utils"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/golang-jwt/jwt"
)

const (
	authPrefix              = "Bearer "
	usernameTenantSeparator = "|"
)

// Structure to define the parser of the Token and what have to be verify in the Token
type JWT struct {
	OrchestrateClaimPath string
	parser               *jwt.Parser
	certificates         []*x509.Certificate
}

func New(cfg *Config) (*JWT, error) {
	return &JWT{
		OrchestrateClaimPath: cfg.OrchestrateClaimPath,
		certificates:         cfg.Certificates,
		parser: &jwt.Parser{
			ValidMethods:         cfg.ValidMethods,
			SkipClaimsValidation: cfg.SkipClaimsValidation,
		},
	}, nil
}

func (checker *JWT) keyFunc(token *jwt.Token) (interface{}, error) {
	for _, cert := range checker.certificates {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); ok {
			return cert.PublicKey.(*rsa.PublicKey), nil
		}
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return cert.PublicKey.(*ecdsa.PublicKey), nil
		}
	}

	return nil, fmt.Errorf("unexpected method: %s", token.Method.Alg())
}

// Parse and verify the validity of the Token (UUID or Access) and return a struct for a JWT (JSON Web Token)
func (checker *JWT) Check(ctx context.Context) (context.Context, error) {
	if len(checker.certificates) == 0 {
		// If no certificate provided we deactivate authentication
		return nil, nil
	}

	if checker == nil || checker.certificates == nil {
		// If no certificate provided we deactivate authentication
		return nil, nil
	}

	// Extract Access Token from context
	bearerToken, ok := authutils.ParseAuth(authPrefix, authutils.AuthorizationFromContext(ctx))
	if !ok {
		return nil, nil
	}

	// Parse and validate token injected in context
	token, err := checker.parser.ParseWithClaims(
		bearerToken,
		&Claims{namespace: checker.OrchestrateClaimPath},
		checker.keyFunc,
	)
	if err != nil {
		return ctx, errors.UnauthorizedError(err.Error())
	}
	if !token.Valid {
		return ctx, errors.UnauthorizedError("invalid Access Token")
	}

	ctx = With(ctx, token)

	// Manage multitenancy
	tenantID, err := multitenancy.TenantID(
		token.Claims.(*Claims).Orchestrate.TenantID,
		multitenancy.TenantIDFromContext(ctx),
	)
	if err != nil {
		return ctx, err
	}

	allowedTenants := multitenancy.AllowedTenants(
		token.Claims.(*Claims).Orchestrate.TenantID,
		multitenancy.TenantIDFromContext(ctx),
	)

	ctx = multitenancy.WithTenantID(ctx, tenantID)
	ctx = multitenancy.WithAllowedTenants(ctx, allowedTenants)

	return ctx, nil
}
