package multitenancy

import (
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
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
func ExtractTenant(checker auth.Checker) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Logger.Tracef("Start handler ExtractTenant")
		if !viper.GetBool(multitenancy.EnabledViperKey) {
			return
		}

		// Extract credentials from envelope metadata
		authorization := txctx.Envelope.GetHeadersValue(AuthorizationMetadata)
		tenantID := txctx.Envelope.GetHeadersValue(TenantIDMetadata)
		authCtx := multitenancy.WithTenantID(authutils.WithAuthorization(txctx.Context(), authorization), tenantID)

		// Control authentication
		checkedCtx, err := checker.Check(authCtx)
		if err != nil {
			e := txctx.AbortWithError(errors.UnauthorizedError(
				err.Error(),
			)).SetComponent(component)
			txctx.Logger.WithError(e).Errorf("Unauthorized: could not authenticate the requester")
			return
		}

		// Attached context enriched with auth information to txctx
		txctx.WithContext(checkedCtx)

		txctx.Logger.Tracef("TenantID extracted and injected in Context with value")
	}
}
