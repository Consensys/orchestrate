package multitenancy

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/auth"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/auth/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
)

const (
	AuthorizationMetadata = "Authorization"
	TenantIDMetadata      = "X-Tenant-ID"
)

// ExtractTenant handler operate:
// - check if the multi-tenancy is enable
// - load /extract the certificate of the Auth Service from config
// - extract the <UUID/Access> Token from the Envelop
// - verify the signature and verify if the certificate from the Token is the same that the loaded certificate => oidc/KeySet
// - extract the Tenant UUID
// - inject the Tenant in the Envelop
func ExtractTenant(multiTenancyEnabled bool, checker auth.Checker) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Logger.Tracef("Start handler ExtractTenant")
		if !multiTenancyEnabled {
			return
		}

		// Extract credentials from envelope metadata
		authorization := txctx.Envelope.GetHeadersValue(AuthorizationMetadata)
		tenantID := txctx.Envelope.GetHeadersValue(TenantIDMetadata)

		if checker == nil {
			ctx := multitenancy.WithTenantID(txctx.Context(), tenantID)
			ctx = multitenancy.WithAllowedTenants(ctx, multitenancy.AllowedTenants(multitenancy.Wildcard, tenantID))
			// Attached context enriched with auth information to txctx
			txctx.WithContext(ctx)
			txctx.Logger.Tracef("TenantID extracted and injected in Context with value")
			return
		}

		ctx := authutils.WithAuthorization(txctx.Context(), authorization)
		ctx = multitenancy.WithTenantID(ctx, tenantID)
		// Control authentication
		checkedCtx, err := checker.Check(ctx)
		if err != nil {
			e := txctx.AbortWithError(errors.UnauthorizedError(
				err.Error(),
			)).SetComponent(component)
			txctx.Logger.WithError(e).Errorf("Unauthorized: could not authenticate the requester")
			return
		}

		// Attached context enriched with auth information to txctx
		txctx.WithContext(checkedCtx)
		txctx.Logger.Tracef("Valid Authentication. TenantID extracted and injected in Context with value")
	}
}
