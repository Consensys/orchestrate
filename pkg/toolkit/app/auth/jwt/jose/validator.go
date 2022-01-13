package jose

import (
	"context"
	"fmt"
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

	opts := []validator.Option{}
	if cfg.OrchestrateClaims != "" {
		opts = append(opts, validator.WithCustomClaims(NewCustomClaims(cfg.OrchestrateClaims)))
	}

	v, err := validator.New(
		jwks.NewCachingProvider(issuerURL, cfg.CacheTTL).KeyFunc,
		validator.RS256,
		issuerURL.String(),
		cfg.Audience,
		opts...,
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
	if claims.CustomClaims != nil {
		if orchestrateUserClaims := claims.CustomClaims.(*CustomClaims).UserClaims; orchestrateUserClaims != nil {
			return orchestrateUserClaims, nil
		}

		return nil, fmt.Errorf("expected custom claims not found")
	}

	// The tenant ID is the "sub" field, then is "tenant_id:username" or "tenant_id"
	claim := &entities.UserClaims{}
	sub := claims.RegisteredClaims.Subject
	pieces := strings.Split(sub, utils.AuthSeparator)
	if len(pieces) == 0 {
		claim.TenantID = pieces[0]
	} else {
		claim.Username = pieces[len(pieces)-1]
		claim.TenantID = strings.Replace(sub, utils.AuthSeparator+claim.Username, "", 1)
	}

	return claim, nil
}
