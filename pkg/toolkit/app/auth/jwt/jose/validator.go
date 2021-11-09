package jose

import (
	"context"
	"net/url"
	"time"

	"github.com/consensys/orchestrate/pkg/types/entities"

	"github.com/auth0/go-jwt-middleware/validate/josev2"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type Validator struct {
	validator *josev2.Validator
}

func NewValidator(cfg *Config) (*Validator, error) {
	issuerURL, err := url.Parse(cfg.IssuerURL)
	if err != nil {
		return nil, err
	}

	expectedClaims := jwt.Expected{Time: time.Now()}
	if len(cfg.Audience) == 0 {
		expectedClaims.Audience = cfg.Audience
	}

	validator, err := josev2.New(
		josev2.NewCachingJWKSProvider(*issuerURL, cfg.CacheTTL).KeyFunc,
		jose.RS256,
		josev2.WithCustomClaims(func() josev2.CustomClaims { return &CustomClaims{} }),
		josev2.WithExpectedClaims(func() jwt.Expected {
			return expectedClaims.WithTime(time.Now())
		}),
	)
	if err != nil {
		return nil, err
	}

	return &Validator{validator: validator}, nil
}

func (v *Validator) ValidateToken(ctx context.Context, token string) (*entities.UserClaims, error) {
	userCtx, err := v.validator.ValidateToken(ctx, token)
	if err != nil {
		// There is no fine-grained handling of the error provided from the package
		return nil, err
	}

	// The tenant ID is the "sub" field, if the "tenant_id" custom field is specified, then tenant ID is "tenant_id"
	tenantID := userCtx.(*josev2.UserContext).Claims.Subject
	tenantIDCustomClaim := userCtx.(*josev2.UserContext).CustomClaims.(*CustomClaims).TenantID
	if tenantIDCustomClaim != "" {
		tenantID = tenantIDCustomClaim
	}

	return &entities.UserClaims{TenantID: tenantID}, nil
}
