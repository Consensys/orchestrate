package jose

import (
	"context"
	"net/url"
	"strings"

	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/consensys/orchestrate/pkg/toolkit/app/auth/utils"
	"github.com/consensys/orchestrate/pkg/types/entities"

	"github.com/auth0/go-jwt-middleware/v2/validator"
)

type Validator struct {
	validator *validator.Validator
}

func NewValidator(cfg *Config) (*Validator, error) {
	issuerURL, err := url.Parse(cfg.IssuerURL)
	if err != nil {
		return nil, err
	}

	v, err := validator.New(
		jwks.NewCachingProvider(issuerURL, cfg.CacheTTL).KeyFunc,
		validator.RS256,
		issuerURL.String(),
		cfg.Audience,
		validator.WithCustomClaims(NewCustomClaims(cfg.OrchestrateClaims)),
	)
	if err != nil {
		return nil, err
	}

	return &Validator{validator: v}, nil
}

func (v *Validator) ValidateToken(ctx context.Context, token string) (*entities.UserClaims, error) {
	userCtx, err := v.validator.ValidateToken(ctx, token)
	if err != nil {
		// There is no fine-grained handling of the error provided from the package
		return nil, err
	}

	claims := userCtx.(*validator.ValidatedClaims)

	if orchestrateUserClaims := claims.CustomClaims.(*CustomClaims).UserClaims; orchestrateUserClaims != nil {
		return orchestrateUserClaims, nil
	}

	// The tenant ID is the "sub" field, then is "tenant_id:username" or "tenant_id"
	sub := claims.RegisteredClaims.Subject
	pieces := strings.Split(sub, utils.AuthSeparator)
	claim := &entities.UserClaims{}
	if len(pieces) == 0 {
		claim.TenantID = pieces[0]
	} else {
		claim.Username = pieces[len(pieces)-1]
		claim.TenantID = strings.Replace(sub, utils.AuthSeparator+claim.Username, "", 1)
	}

	return claim, nil
}
