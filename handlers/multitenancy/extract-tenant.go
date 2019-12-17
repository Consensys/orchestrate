package multitenancy

import (
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication"
)

// ExtractTenant handler operate:
// - check if the multi-tenancy is enable
// - load /extract the certificate of the Auth Service from config
// - extract the <ID/Access> Token from the Envelop
// - verify the signature and verify if the certificate from the Token is the same that the loaded certificate => oidc/KeySet
// - extract the Tenant ID
// - inject the Tenant in the Envelop
func ExtractTenant(m authentication.Manager) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		if !viper.GetBool(multitenancy.EnabledViperKey) {
			// Run the next middlewares
			txctx.Next()
		}

		rawToken, err := m.Extract(txctx.Envelope)
		if err != nil {
			e := txctx.AbortWithError(errors.UnauthorizedError(
				"Token Not Found: " + err.Error(),
			)).SetComponent(component)
			txctx.Logger.WithError(e).Errorf("Token Not Found: could extract the ID / Access Token from the envelop")
			return
		}

		token, err := m.Verify(rawToken)
		if err != nil {
			e := txctx.AbortWithError(errors.UnauthorizedError(
				err.Error(),
			)).SetComponent(component)
			txctx.Logger.WithError(e).Errorf("Unauthorized: could not authenticate the requester")
			return
		}

		tenantPath := viper.GetString(authentication.TenantNamespaceViperKey)

		tenantIDValue, ok := token.Claims.(jwt.MapClaims)[tenantPath+authentication.TenantIDKey].(string)
		if !ok {
			err := fmt.Errorf("not able to retrieve the tenant ID: The tenant_id is not present in the ID / Access Token")
			_ = txctx.AbortWithError(errors.NotFoundError(
				err.Error(),
			)).SetComponent(component)
			txctx.Logger.Error(err)
			return
		}
		// Add the Token information and the Tenant Id in the go Context into the transaction Context
		txctx.Set(authentication.TokenInfoKey, token)
		txctx.Set(authentication.TenantIDKey, tenantIDValue)

		txctx.Next()
	}
}
