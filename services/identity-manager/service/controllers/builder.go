package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/identity-manager/use-cases"
)

//go:generate swag init --dir . --generalInfo builder.go --output ../../../../public/swagger-specs/services/identity-manager
//go:generate rm ../../../../public/swagger-specs/services/identity-manager/docs.go ../../../../public/swagger-specs/services/identity-manager/swagger.yaml

// @title Identity Manager API
// @version 2.0
// @description PegaSys Orchestrate Identity API. Enables dynamic management of identities.
// @description Identities correspond to an Ethereum accounts. It can be a user account or a deployed smart contract. By usage of the generated cryptographic key pair, identities can be used to sign/verify and to encrypt/decrypt messages.

// @contact.name Contact PegaSys Orchestrate
// @contact.url https://pegasys.tech/contact/
// @contact.email support@pegasys.tech

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key

// @securityDefinitions.apikey JWTAuth
// @in header
// @name Authorization

type Builder struct {
	identitiesCtrl *IdentitiesController
}

func NewBuilder(ucs usecases.IdentityUseCases) *Builder {
	return &Builder{
		identitiesCtrl: NewIdentitiesController(ucs),
	}
}

func (b *Builder) Build(ctx context.Context, _ string, configuration interface{}, _ func(response *http.Response) error) (http.Handler, error) {
	cfg, ok := configuration.(*dynamic.Identity)
	if !ok {
		return nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	router := mux.NewRouter()
	b.identitiesCtrl.Append(router)

	return router, nil
}
