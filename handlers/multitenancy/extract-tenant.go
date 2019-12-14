package multitenancy

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication"
)

const component = "handler.multitenancy"
const TenantIDKey = "tenant_id"

// ExtractTenant handler operate:
// - check if the multi-tenancy is enable
// - load /extract the certificate of the Auth Service from config
// - extract the <ID/Access> Token from the Envelop
// - verify the signature and verify if the certificate from the Token is the same that the loaded certificate => oidc/KeySet
// - extract the Tenant ID
// - inject the Tenant in the Envelop
func ExtractTenant(m authentication.Manager) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		if !viper.GetBool(EnabledViperKey) {
			// Run the next middlewares
			txctx.Next()
		}

		rawToken, err := m.Extract(txctx.Envelope)
		if err != nil {
			e := txctx.AbortWithError(errors.NotFoundError(
				err.Error(),
			)).SetComponent(component)
			txctx.Logger.WithError(e).Errorf("Token Not Found: could extract the ID / Access Token from the envelop")
			return
		}

		token, err := m.Verify(rawToken)
		if err != nil {
			e := txctx.AbortWithError(errors.UnauthenticatedError(
				err.Error(),
			)).SetComponent(component)
			txctx.Logger.WithError(e).Errorf("Unauthenticated: could not authenticate the requester")
			return
		}

		tenantPath := viper.GetString(TenantNamespaceViperKey)

		tenantIDValue, ok := token.Claims.(jwt.MapClaims)[tenantPath+TenantIDKey].(string)
		if !ok {
			err := fmt.Errorf("not able to retrieve the tenant ID: The tenant_id is not present in the ID / Access Token")
			_ = txctx.AbortWithError(errors.NotFoundError(
				err.Error(),
			)).SetComponent(component)
			txctx.Logger.Error(err)
			return
		}

		txctx.Set(TenantIDKey, tenantIDValue)

		txctx.Next()
	}
}
