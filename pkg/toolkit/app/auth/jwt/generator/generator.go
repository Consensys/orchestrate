package generator

import (
	"context"
	"crypto/rsa"
	"time"

	"github.com/consensys/orchestrate/pkg/toolkit/app/auth/jwt"
	"github.com/consensys/orchestrate/pkg/toolkit/tls/certificate"
	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/handler/openid"
	fositejwt "github.com/ory/fosite/token/jwt"
)

type JWTGenerator struct {
	ClaimsNamespace string
	privateKey      *rsa.PrivateKey
}

func New(cfg *Config) (*JWTGenerator, error) {
	cert, err := certificate.X509(cfg.KeyPair)
	if err != nil {
		return nil, err
	}

	return &JWTGenerator{
		ClaimsNamespace: cfg.ClaimsNamespace,
		privateKey:      cert.PrivateKey.(*rsa.PrivateKey),
	}, nil
}

func (j *JWTGenerator) GenerateAccessTokenWithTenantID(tenantID string, ttl time.Duration) (string, error) {
	customClaims := map[string]interface{}{
		j.ClaimsNamespace: &jwt.OrchestrateClaims{
			TenantID: tenantID,
		}}
	return j.GenerateAccessToken(customClaims, ttl)
}

func (j *JWTGenerator) GenerateAccessToken(customClaims map[string]interface{}, ttl time.Duration) (tokenValue string, err error) {
	jwtGenerator := &oauth2.DefaultJWTStrategy{
		JWTStrategy: &fositejwt.RS256JWTStrategy{
			PrivateKey: j.privateKey,
		},
	}
	tokenRequest := &fosite.Request{
		GrantedAudience: []string{"https://auth0.com/api/v2/"},
		GrantedScope:    []string{"read:users", "update:users", "create:users"},
		Client: &fosite.DefaultClient{
			ID:     "App-test",
			Secret: []byte("mysecret"),
		},
		Session: &oauth2.JWTSession{
			JWTClaims: &fositejwt.JWTClaims{
				Issuer:    "Orchestrate",
				Subject:   "e2e-test",
				IssuedAt:  time.Now().UTC(),
				NotBefore: time.Now().UTC(),
				Extra:     customClaims,
			},
			JWTHeader: &fositejwt.Headers{
				Extra: make(map[string]interface{}),
			},
			ExpiresAt: map[fosite.TokenType]time.Time{
				fosite.AccessToken: time.Now().UTC().Add(ttl),
			},
		},
	}
	// The access Token contain already the signature
	accessToken, _, err := jwtGenerator.GenerateAccessToken(context.Background(), tokenRequest)
	return accessToken, err
}

func (j *JWTGenerator) GenerateIDToken(customClaims map[string]interface{}) (tokenValue string, err error) {

	jwtGenerator := &openid.DefaultStrategy{
		JWTStrategy: &fositejwt.RS256JWTStrategy{
			PrivateKey: j.privateKey,
		},
	}

	tokenRequest := fosite.NewAccessRequest(&openid.DefaultSession{
		Claims: &fositejwt.IDTokenClaims{
			Issuer:   "Orchestrate",
			Subject:  "e2e-test",
			Audience: []string{"https://auth0.com/api/v2/"},
			IssuedAt: time.Now().UTC(),
			Extra:    customClaims,
		},
		Headers: &fositejwt.Headers{},
	})

	idToken, err := jwtGenerator.GenerateIDToken(context.Background(), tokenRequest)

	return idToken, err
}
