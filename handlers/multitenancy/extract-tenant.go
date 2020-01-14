package multitenancy

import (
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/token"

	"github.com/dgrijalva/jwt-go"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"

	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
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
		txctx.Logger.Tracef("Start handler ExtractTenant")
		if !viper.GetBool(multitenancy.EnabledViperKey) {
			// Run the next middlewares
			txctx.Next()
		} else {
			rawToken, err := token.Extract(txctx.Envelope)
			if err != nil {
				e := txctx.AbortWithError(errors.UnauthorizedError(
					"Token Not Found: " + err.Error(),
				)).SetComponent(component)
				txctx.Logger.WithError(e).Errorf("Token Not Found: could extract the ID / Access Token from the envelop")
				return
			}

			tokenStruct, err := m.Verify(rawToken)
			if err != nil {
				e := txctx.AbortWithError(errors.UnauthorizedError(
					err.Error(),
				)).SetComponent(component)
				txctx.Logger.WithError(e).Errorf("Unauthorized: could not authenticate the requester")
				return
			}

			tenantPath := viper.GetString(authentication.TenantNamespaceViperKey)

			tenantIDValue, ok := tokenStruct.Claims.(jwt.MapClaims)[tenantPath+authentication.TenantIDKey].(string)
			if !ok {
				err := fmt.Errorf("not able to retrieve the tenant ID: The tenant_id is not present in the ID / Access Token")
				_ = txctx.AbortWithError(errors.NotFoundError(
					err.Error(),
				)).SetComponent(component)
				txctx.Logger.Error(err)
				return
			}
			// Add the Token information and the Tenant Id in the go Context into the transaction Context
			txctx.Set(authentication.TokenInfoKey, tokenStruct)
			txctx.Set(authentication.TokenRawKey, rawToken)
			txctx.Set(authentication.TenantIDKey, tenantIDValue)

			txctx.Logger.Tracef("TenantID extracted and injected in Context with value: %s", tenantIDValue)

			txctx.Next()
		}
	}
}
