package multitenancy

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"

	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication"
)

const AuthorizationMetadata = "Authorization"

// ExtractTenant handler operate:
// - check if the multi-tenancy is enable
// - load /extract the certificate of the Auth Service from config
// - extract the <ID/Access> Token from the Envelop
// - verify the signature and verify if the certificate from the Token is the same that the loaded certificate => oidc/KeySet
// - extract the Tenant ID
// - inject the Tenant in the Envelop
func ExtractTenant(auth authentication.Auth) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Logger.Tracef("Start handler ExtractTenant")
		if !viper.GetBool(multitenancy.EnabledViperKey) {
			return
		}

		// Extract credentials from envelope metadata
		authorization, _ := txctx.Envelope.GetMetadataValue(AuthorizationMetadata)

		// Control authentication
		checkedCtx, err := auth.Check(authutils.WithAuthorization(txctx.Context(), authorization))
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
