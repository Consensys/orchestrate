package jwt

import (
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/multitenancy"
	authutils "github.com/ConsenSys/orchestrate/pkg/toolkit/app/auth/utils"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/tls/certificate"
	"github.com/dgrijalva/jwt-go"
)

// Structure to define the parser of the Token and what have to be verify in the Token
type JWT struct {
	cfg    *Config
	parser *jwt.Parser
	cert   *x509.Certificate
}

func New(cfg *Config) (*JWT, error) {
	cert, err := certificate.X509KeyPair(cfg.Certificate, nil)
	if err != nil {
		return nil, err
	}

	return &JWT{
		cfg:  cfg,
		cert: cert.Leaf,
		parser: &jwt.Parser{
			ValidMethods:         cfg.ValidMethods,
			SkipClaimsValidation: cfg.SkipClaimsValidation,
		},
	}, nil
}

func (checker *JWT) key(token *jwt.Token) (interface{}, error) {
	switch token.Method.Alg() {
	case "RS256", "RS384", "RS512":
		pubKey, ok := checker.cert.PublicKey.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("jwt: certificate is not an RSA public key")
		}
		return pubKey, nil
	case "ES256", "ES384", "ES512":
		pubKey, ok := checker.cert.PublicKey.(*ecdsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("jwt: certificate is not an ECDSA public key")
		}
		return pubKey, nil
	default:
		return nil, fmt.Errorf("jwt: unsupported token method signature %q", token.Method.Alg())
	}
}

const authPrefix = "Bearer "

// Parse and verify the validity of the Token (UUID or Access) and return a struct for a JWT (JSON Web Token)
func (checker *JWT) Check(ctx context.Context) (context.Context, error) {
	if checker == nil || checker.cert == nil {
		// If no certificate provided we deactivate authentication
		return nil, nil
	}

	// Extract Access Token from context
	bearer, ok := authutils.ParseAuth(authPrefix, authutils.AuthorizationFromContext(ctx))
	if !ok {
		return nil, nil
	}

	// Parse and validate token injected in context
	token, err := checker.parser.ParseWithClaims(
		bearer,
		&Claims{namespace: checker.cfg.ClaimsNamespace},
		checker.key,
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
