package jose

import (
	"context"
	"net/url"
	"strings"
	"time"

	authutils "github.com/consensys/orchestrate/pkg/toolkit/app/auth/utils"
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
		josev2.WithCustomClaims(func() josev2.CustomClaims { return NewCustomClaims(cfg.OrchestrateClaims) }),
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

	if orchestrateUserClaims := userCtx.(*josev2.UserContext).CustomClaims.(*CustomClaims).UserClaims; orchestrateUserClaims != nil {
		return orchestrateUserClaims, nil
	}

	// The tenant ID is the "sub" field, then is "tenant_id:username" or "tenant_id"
	sub := userCtx.(*josev2.UserContext).Claims.Subject
	pieces := strings.Split(sub, authutils.AuthSeparator)
	claim := &entities.UserClaims{}
	if len(pieces) == 0 {
		claim.TenantID = pieces[0]
	} else {
		claim.Username = pieces[len(pieces)-1]
		claim.TenantID = strings.Replace(sub, authutils.AuthSeparator+claim.Username, "", 1)
	}

	return claim, nil
}
